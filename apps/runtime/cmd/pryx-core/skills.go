package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"pryx-core/internal/config"
	"pryx-core/internal/skills"
)

func runSkills(args []string) int {
	if len(args) < 1 {
		skillsUsage()
		return 2
	}

	cmd := args[0]
	cfg := config.Load()

	switch cmd {
	case "list", "ls":
		return runListSkills(args[1:], cfg)
	case "info":
		return runInfoSkill(args[1:], cfg)
	case "check":
		return runCheckSkills(args[1:], cfg)
	case "enable":
		return runEnableSkill(args[1:], cfg)
	case "disable":
		return runDisableSkill(args[1:], cfg)
	case "install":
		return runInstallSkill(args[1:], cfg)
	case "uninstall":
		return runUninstallSkill(args[1:], cfg)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		skillsUsage()
		return 2
	}
}

func runListSkills(args []string, cfg *config.Config) int {
	eligibleOnly := false
	jsonOutput := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--eligible", "-e":
			eligibleOnly = true
		case "--json", "-j":
			jsonOutput = true
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			skillsUsage()
			return 2
		}
	}

	opts := skills.DefaultOptions()

	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover skills: %v\n", err)
		return 1
	}

	var skillsToDisplay []skills.Skill
	for _, skill := range skillsRepo.List() {
		if eligibleOnly {
			if !skill.Eligible {
				continue
			}
		}
		skillsToDisplay = append(skillsToDisplay, skill)
	}

	// Sort by name
	sort.Slice(skillsToDisplay, func(i, j int) bool {
		return skillsToDisplay[i].ID < skillsToDisplay[j].ID
	})

	if jsonOutput {
		data, err := json.MarshalIndent(skillsToDisplay, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal skills: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
	} else {
		fmt.Printf("Available Skills (%d)\n", len(skillsToDisplay))
		fmt.Println(strings.Repeat("=", 51))
		for _, skill := range skillsToDisplay {
			status := ""
			if !skill.Enabled {
				status = " (disabled)"
			} else if skill.Eligible {
				status = " ✓"
			} else {
				status = " ⚠"
			}

			title := skill.Frontmatter.Name
			if title == "" {
				title = skill.ID
			}
			fmt.Printf("%s %s: %s\n", status, skill.ID, title)
			if skill.Frontmatter.Description != "" {
				fmt.Printf("  %s\n", skill.Frontmatter.Description)
			}
			fmt.Printf("  Source: %s, Enabled: %v\n", skill.Source, skill.Enabled)
		}
	}

	return 0
}

func runInfoSkill(args []string, cfg *config.Config) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: skill name required\n")
		return 2
	}
	name := args[0]

	opts := skills.DefaultOptions()

	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover skills: %v\n", err)
		return 1
	}

	skill, found := skillsRepo.Get(name)
	if !found {
		fmt.Fprintf(os.Stderr, "Error: skill not found: %s\n", name)
		return 1
	}

	fmt.Printf("Skill: %s\n", skill.ID)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Title:       %s\n", skill.Frontmatter.Name)
	fmt.Printf("Description: %s\n", skill.Frontmatter.Description)
	fmt.Printf("Source:      %s\n", skill.Source)
	fmt.Printf("Path:        %s\n", skill.Path)
	fmt.Printf("Enabled:     %v\n", skill.Enabled)
	fmt.Printf("Eligible:    %v\n", skill.Eligible)

	if len(skill.Frontmatter.Metadata.Pryx.Requires.Bins) > 0 {
		fmt.Printf("Required binaries: %s\n", strings.Join(skill.Frontmatter.Metadata.Pryx.Requires.Bins, ", "))
	}
	if len(skill.Frontmatter.Metadata.Pryx.Requires.Env) > 0 {
		fmt.Printf("Required env vars: %s\n", strings.Join(skill.Frontmatter.Metadata.Pryx.Requires.Env, ", "))
	}
	if len(skill.Frontmatter.Metadata.Pryx.Install) > 0 {
		fmt.Printf("Installers: %d\n", len(skill.Frontmatter.Metadata.Pryx.Install))
		for i, installer := range skill.Frontmatter.Metadata.Pryx.Install {
			fmt.Printf("  [%d] %s %s\n", i+1, installer.Command, strings.Join(installer.Args, " "))
		}
	}

	return 0
}

func runCheckSkills(args []string, cfg *config.Config) int {
	opts := skills.DefaultOptions()

	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover skills: %v\n", err)
		return 1
	}

	fmt.Printf("Skills Check\n")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	allSkills := skillsRepo.List()
	if len(allSkills) == 0 {
		fmt.Println("No skills found.")
		return 0
	}

	// Check skills
	validCount := 0
	invalidCount := 0
	issues := 0

	for _, skill := range allSkills {
		issuesInSkill := 0

		// Check SKILL.md exists
		if skill.Path == "" {
			fmt.Printf("✗ %s: No path defined\n", skill.ID)
			issuesInSkill++
			issues++
			continue
		}

		// Check required fields
		if skill.ID == "" {
			fmt.Printf("✗ %s: Missing name\n", skill.ID)
			issuesInSkill++
			issues++
		}
		// Version field doesn't exist in current Frontmatter
		if skill.Frontmatter.Description == "" {
			fmt.Printf("✗ %s: Missing description\n", skill.ID)
			issuesInSkill++
			issues++
		}

		// Check prompts
		if len(skill.SystemPrompt) == 0 {
			body, _ := skill.Body()
			if strings.TrimSpace(body) == "" {
				fmt.Printf("⚠ %s: Empty system prompt\n", skill.ID)
				issuesInSkill++
				issues++
			}
		}

		if issuesInSkill == 0 {
			fmt.Printf("✓ %s: All checks passed\n", skill.ID)
			validCount++
		} else {
			invalidCount++
		}
	}

	fmt.Println()
	fmt.Printf("Summary:\n")
	fmt.Printf("  Total Skills:  %d\n", len(allSkills))
	fmt.Printf("  Valid Skills:  %d\n", validCount)
	fmt.Printf("  Invalid Skills: %d\n", invalidCount)
	fmt.Printf("  Total Issues:  %d\n", issues)

	if issues == 0 {
		fmt.Println()
		fmt.Printf("✓ All skills are properly configured\n")
		return 0
	} else {
		fmt.Println()
		fmt.Printf("✗ Found %d issues across %d skills\n", issues, invalidCount)
		return 1
	}
}

func runEnableSkill(args []string, cfg *config.Config) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: skill name required\n")
		return 2
	}
	name := args[0]

	opts := skills.DefaultOptions()

	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover skills: %v\n", err)
		return 1
	}

	_, found := skillsRepo.Get(name)
	if !found {
		fmt.Fprintf(os.Stderr, "Error: skill not found: %s\n", name)
		return 1
	}

	configPath := skills.EnabledConfigPath()
	enabledCfg, err := skills.LoadEnabledConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load skills config: %v\n", err)
		return 1
	}

	if enabledCfg.EnabledSkills[name] {
		fmt.Printf("ℹ Skill %s is already enabled\n", name)
	} else {
		enabledCfg.EnabledSkills[name] = true
		if err := skills.SaveEnabledConfig(configPath, enabledCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to save skills config: %v\n", err)
			return 1
		}
		fmt.Printf("✓ Enabled skill: %s\n", name)
	}

	return 0
}

func runDisableSkill(args []string, cfg *config.Config) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: skill name required\n")
		return 2
	}
	name := args[0]

	opts := skills.DefaultOptions()

	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to discover skills: %v\n", err)
		return 1
	}

	_, found := skillsRepo.Get(name)
	if !found {
		fmt.Fprintf(os.Stderr, "Error: skill not found: %s\n", name)
		return 1
	}

	configPath := skills.EnabledConfigPath()
	enabledCfg, err := skills.LoadEnabledConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load skills config: %v\n", err)
		return 1
	}

	if !enabledCfg.EnabledSkills[name] {
		fmt.Printf("ℹ Skill %s is already disabled\n", name)
	} else {
		delete(enabledCfg.EnabledSkills, name)
		if err := skills.SaveEnabledConfig(configPath, enabledCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to save skills config: %v\n", err)
			return 1
		}
		fmt.Printf("✓ Disabled skill: %s\n", name)
	}

	return 0
}

func runInstallSkill(args []string, cfg *config.Config) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: skill name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pryx-core skills install <name> [--from <path|url>]\n")
		return 2
	}
	name := args[0]

	opts := skills.DefaultOptions()

	skill, err := installSkillFromSource(name, args[1:], opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to install skill: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Skill installed successfully: %s\n", skill.ID)
	fmt.Printf("  Name: %s\n", skill.Frontmatter.Name)
	fmt.Printf("  Path: %s\n", skill.Path)
	fmt.Printf("  Source: %s\n", skill.Source)

	if len(skill.Frontmatter.Metadata.Pryx.Install) > 0 {
		fmt.Printf("\n  Running %d installer(s)...\n", len(skill.Frontmatter.Metadata.Pryx.Install))
		if err := runSkillInstallers(skill); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: installation steps failed: %v\n", err)
			fmt.Println("  Skill files installed but dependencies may need manual setup.")
		} else {
			fmt.Println("  ✓ All installation steps completed")
		}
	}

	fmt.Printf("\nEnable the skill with: pryx-core skills enable %s\n", skill.ID)
	return 0
}

func installSkillFromSource(name string, args []string, opts skills.Options) (*skills.Skill, error) {
	from := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--from" && i+1 < len(args) {
			from = args[i+1]
			i++
		}
	}

	skillPath := filepath.Join(opts.ManagedRoot, name)

	if from == "" {
		return nil, fmt.Errorf("no source specified. Use --from <path|url|bundled/skill-name>")
	}

	if strings.HasPrefix(from, "bundled/") {
		bundledName := strings.TrimPrefix(from, "bundled/")
		bundledPath := filepath.Join(opts.BundledRoot, bundledName)
		if _, err := os.Stat(bundledPath); err != nil {
			return nil, fmt.Errorf("bundled skill not found: %s", bundledName)
		}
		from = bundledPath
	}

	if strings.HasPrefix(from, "http://") || strings.HasPrefix(from, "https://") {
		return installSkillFromURL(name, from, skillPath)
	}

	if _, err := os.Stat(from); err != nil {
		return nil, fmt.Errorf("source path not found: %s", from)
	}

	return installSkillFromPath(name, from, skillPath)
}

func installSkillFromPath(name, sourcePath, targetPath string) (*skills.Skill, error) {
	sourceSkillPath := filepath.Join(sourcePath, "SKILL.md")
	if _, err := os.Stat(sourceSkillPath); err != nil {
		sourceSkillPath = sourcePath
		if !strings.HasSuffix(sourceSkillPath, ".md") {
			return nil, fmt.Errorf("source does not contain SKILL.md: %s", sourcePath)
		}
	}

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create skill directory: %w", err)
	}

	targetSkillPath := filepath.Join(targetPath, "SKILL.md")

	data, err := os.ReadFile(sourceSkillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source skill: %w", err)
	}

	if err := os.WriteFile(targetSkillPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write skill file: %w", err)
	}

	skill := skills.Skill{
		ID:     name,
		Source: skills.SourceManaged,
		Path:   targetPath,
	}

	return &skill, nil
}

func installSkillFromURL(name, url, targetPath string) (*skills.Skill, error) {
	return nil, fmt.Errorf("URL-based installation not yet implemented")
}

func runSkillInstallers(skill *skills.Skill) error {
	if len(skill.Frontmatter.Metadata.Pryx.Install) == 0 {
		return nil
	}

	for i, installer := range skill.Frontmatter.Metadata.Pryx.Install {
		fmt.Printf("  [%d/%d] Running: %s %s\n", i+1, len(skill.Frontmatter.Metadata.Pryx.Install), installer.Command, strings.Join(installer.Args, " "))

		cmd := exec.Command(installer.Command, installer.Args...)
		cmd.Dir = skill.Path
		if len(installer.Env) > 0 {
			cmd.Env = append(os.Environ(), installer.Env...)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("installer %d failed: %w\nOutput: %s", i+1, err, string(output))
		}

		if len(output) > 0 {
			fmt.Printf("    %s\n", string(output))
		}
	}

	return nil
}

func runUninstallSkill(args []string, cfg *config.Config) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: skill name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pryx-core skills uninstall <name> [--force]\n")
		return 2
	}
	name := args[0]

	force := false
	for _, arg := range args[1:] {
		if arg == "--force" || arg == "-f" {
			force = true
		}
	}

	opts := skills.DefaultOptions()
	skillPath := filepath.Join(opts.ManagedRoot, name)

	opts = skills.DefaultOptions()
	skillsRepo, err := skills.Discover(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to discover skills: %v\n", err)
	}

	if skill, found := skillsRepo.Get(name); found {
		if skill.Source != skills.SourceManaged {
			if !force {
				fmt.Fprintf(os.Stderr, "Error: skill '%s' is %s (not managed). Use --force to remove.\n", name, skill.Source)
				return 1
			}
			fmt.Printf("Warning: removing %s skill '%s' (forced)\n", skill.Source, name)
		}

		if skill.Enabled {
			fmt.Printf("Disabling skill '%s'...\n", name)
			enabledCfg, _ := skills.LoadEnabledConfig(skills.EnabledConfigPath())
			delete(enabledCfg.EnabledSkills, name)
			skills.SaveEnabledConfig(skills.EnabledConfigPath(), enabledCfg)
		}
	}

	if _, err := os.Stat(skillPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: skill not found at %s\n", skillPath)
		return 1
	}

	fmt.Printf("Removing skill directory: %s\n", skillPath)
	if err := os.RemoveAll(skillPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to remove skill: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Skill uninstalled: %s\n", name)
	return 0
}

func skillsUsage() {
	fmt.Println("pryx-core skills - Manage Pryx skills")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list [--eligible] [--json]        List available skills")
	fmt.Println("  info <name>                       Show skill details")
	fmt.Println("  check                             Check all skills for issues")
	fmt.Println("  enable <name>                     Enable a skill")
	fmt.Println("  disable <name>                    Disable a skill")
	fmt.Println("  install <name>                    Install a skill")
	fmt.Println("  uninstall <name>                  Uninstall a skill")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --eligible, -e                    Show only eligible skills")
	fmt.Println("  --json, -j                        Output in JSON format")
}
