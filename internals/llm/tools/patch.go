package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/omnitrix-sh/cli/internals/config"
	"github.com/omnitrix-sh/cli/internals/diff"
	"github.com/omnitrix-sh/cli/internals/history"
	"github.com/omnitrix-sh/cli/internals/lsp"
)

type PatchParams struct {
	PatchText string `json:"patch_text"`
}

type PatchResponseMetadata struct {
	FilesChanged []string `json:"files_changed"`
	Additions    int      `json:"additions"`
	Removals     int      `json:"removals"`
}

type patchTool struct {
	lspClients map[string]*lsp.Client
	files      history.Service
}

const (
	PatchToolName    = "patch"
	patchDescription = `Applies a patch to multiple files in one operation. This tool is useful for making coordinated changes across multiple files.

The patch text must follow this format:
*** Begin Patch
*** Update File: /path/to/file
@@ Context line (unique within the file)
 Line to keep
-Line to remove
+Line to add
 Line to keep
*** Add File: /path/to/new/file
+Content of the new file
+More content
*** Delete File: /path/to/file/to/delete
*** End Patch

Before using this tool:
1. Use the FileRead tool to understand the files' contents and context
2. Verify all file paths are correct (use the LS tool)

CRITICAL REQUIREMENTS FOR USING THIS TOOL:

1. UNIQUENESS: Context lines MUST uniquely identify the specific sections you want to change
2. PRECISION: All whitespace, indentation, and surrounding code must match exactly
3. VALIDATION: Ensure edits result in idiomatic, correct code
4. PATHS: Always use absolute file paths (starting with /)

The tool will apply all changes in a single atomic operation.`
)

func NewPatchTool(lspClients map[string]*lsp.Client, files history.Service) BaseTool {
	return &patchTool{
		lspClients: lspClients,
		files:      files,
	}
}

func (p *patchTool) Info() ToolInfo {
	return ToolInfo{
		Name:        PatchToolName,
		Description: patchDescription,
		Parameters: map[string]any{
			"patch_text": map[string]any{
				"type":        "string",
				"description": "The full patch text that describes all changes to be made",
			},
		},
		Required: []string{"patch_text"},
	}
}

func (p *patchTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params PatchParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse("invalid parameters"), nil
	}

	if params.PatchText == "" {
		return NewTextErrorResponse("patch_text is required"), nil
	}

	filesToRead := diff.IdentifyFilesNeeded(params.PatchText)
	for _, filePath := range filesToRead {
		absPath := filePath
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}

		if getLastReadTime(absPath).IsZero() {
			return NewTextErrorResponse(fmt.Sprintf("you must read the file %s before patching it. Use the FileRead tool first", filePath)), nil
		}

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				return NewTextErrorResponse(fmt.Sprintf("file not found: %s", absPath)), nil
			}
			return ToolResponse{}, fmt.Errorf("failed to access file: %w", err)
		}

		if fileInfo.IsDir() {
			return NewTextErrorResponse(fmt.Sprintf("path is a directory, not a file: %s", absPath)), nil
		}

		modTime := fileInfo.ModTime()
		lastRead := getLastReadTime(absPath)
		if modTime.After(lastRead) {
			return NewTextErrorResponse(
				fmt.Sprintf("file %s has been modified since it was last read (mod time: %s, last read: %s)",
					absPath, modTime.Format(time.RFC3339), lastRead.Format(time.RFC3339),
				)), nil
		}
	}

	filesToAdd := diff.IdentifyFilesAdded(params.PatchText)
	for _, filePath := range filesToAdd {
		absPath := filePath
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}

		_, err := os.Stat(absPath)
		if err == nil {
			return NewTextErrorResponse(fmt.Sprintf("file already exists and cannot be added: %s", absPath)), nil
		} else if !os.IsNotExist(err) {
			return ToolResponse{}, fmt.Errorf("failed to check file: %w", err)
		}
	}

	currentFiles := make(map[string]string)
	for _, filePath := range filesToRead {
		absPath := filePath
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			return ToolResponse{}, fmt.Errorf("failed to read file %s: %w", absPath, err)
		}
		currentFiles[filePath] = string(content)
	}

	patch, fuzz, err := diff.TextToPatch(params.PatchText, currentFiles)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to parse patch: %s", err)), nil
	}

	if fuzz > 3 {
		return NewTextErrorResponse(fmt.Sprintf("patch contains fuzzy matches (fuzz level: %d). Please make your context lines more precise", fuzz)), nil
	}

	commit, err := diff.PatchToCommit(patch, currentFiles)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to create commit from patch: %s", err)), nil
	}

	sessionID, _ := GetContextValues(ctx)
	if sessionID == "" {
		return ToolResponse{}, fmt.Errorf("session ID is required for creating a patch")
	}

	err = diff.ApplyCommit(commit, func(path string, content string) error {
		absPath := path
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}

		dir := filepath.Dir(absPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", absPath, err)
		}

		return os.WriteFile(absPath, []byte(content), 0o644)
	}, func(path string) error {
		absPath := path
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}
		return os.Remove(absPath)
	})
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("failed to apply patch: %s", err)), nil
	}

	changedFiles := []string{}
	totalAdditions := 0
	totalRemovals := 0

	for path, change := range commit.Changes {
		absPath := path
		if !filepath.IsAbs(absPath) {
			wd := config.WorkingDirectory()
			absPath = filepath.Join(wd, absPath)
		}
		changedFiles = append(changedFiles, absPath)

		oldContent := ""
		if change.OldContent != nil {
			oldContent = *change.OldContent
		}

		newContent := ""
		if change.NewContent != nil {
			newContent = *change.NewContent
		}

		_, additions, removals := diff.GenerateDiff(oldContent, newContent, path)
		totalAdditions += additions
		totalRemovals += removals

		if p.files != nil {
			file, err := p.files.GetByPathAndSession(ctx, absPath, sessionID)
			if err != nil && change.Type != diff.ActionAdd {
				_, _ = p.files.Create(ctx, sessionID, absPath, oldContent)
			}

			if err == nil && change.Type != diff.ActionAdd && file.Content != oldContent {
				_, _ = p.files.CreateVersion(ctx, sessionID, absPath, oldContent)
			}

			if change.Type == diff.ActionDelete {
				_, _ = p.files.CreateVersion(ctx, sessionID, absPath, "")
			} else {
				_, _ = p.files.CreateVersion(ctx, sessionID, absPath, newContent)
			}
		}

		recordFileWrite(absPath)
		recordFileRead(absPath)
	}

	for _, filePath := range changedFiles {
		waitForLspDiagnostics(ctx, filePath, p.lspClients)
	}

	result := fmt.Sprintf("Patch applied successfully. %d files changed, %d additions, %d removals",
		len(changedFiles), totalAdditions, totalRemovals)

	diagnosticsText := ""
	for _, filePath := range changedFiles {
		diagnosticsText += getDiagnostics(filePath, p.lspClients)
	}

	if diagnosticsText != "" {
		result += "\n\nDiagnostics:\n" + diagnosticsText
	}

	return WithResponseMetadata(
		NewTextResponse(result),
		PatchResponseMetadata{
			FilesChanged: changedFiles,
			Additions:    totalAdditions,
			Removals:     totalRemovals,
		}), nil
}
