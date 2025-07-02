package endpoints

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/config"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
)

const (
	caddyfilePath = "/etc/caddy/Caddyfile"
)

type SetupDomainRequest struct {
	Domain             string `json:"domain"`
	MainCategoryNumber int    `json:"mainCategoryNumber"`
}

func (c *SetupDomainRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	// Add validation logic for the SetupDomainRequest struct
	return nil
}

func (s *service) SetupDomain(ctx context.Context, input SetupDomainRequest) response.ErrorResponse {
	// Remove "https://" prefix from the domain
	domain := strings.TrimPrefix(input.Domain, "https://")

	// Check if the domain is already registered
	if domainExists(domain) {
		return response.ErrorBadRequest("دامنه قبلا ثبت شده است")
	}

	// Check if the domain IP matches the server IP
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Printf("Error looking up IP for domain %s: %v", domain, err)
		return response.ErrorNotFound("لطفاً ابتدا تنظیمات IP زیردامنه را انجام دهید.")
	}

	serverIP := net.ParseIP(config.AppConfig.IPServer)
	for _, ip := range ips {
		if !ip.Equal(serverIP) {
			log.Printf("Domain IP %s does not match server IP %s", ip.String(), serverIP.String())
			return response.ErrorBadRequest("آدرس IP دامنه با آدرس سرور مطابقت ندارد")
		}
	}

	// Log the input data
	log.Printf("Received SetupDomainRequest: %+v\n", domain)

	// Implement the logic to update Caddy configuration based on the input
	err = updateCaddyConfigForMainCategory(domain, fmt.Sprint(input.MainCategoryNumber))
	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.ErrorNotFound("خطا در تنظیم دامنه")
	}

	return response.ErrorResponse{}
}

func domainExists(domain string) bool {
	// Open the Caddyfile for reading
	file, err := os.Open(caddyfilePath)
	if err != nil {
		// If there's an error opening the file, consider domain does not exist
		return false
	}
	defer file.Close()

	// Scan through the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains the domain
		if strings.Contains(line, domain) {
			// Domain already exists in Caddyfile
			return true
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		// If there's an error scanning the file, consider domain does not exist
		return false
	}

	// If the domain is not found in the file, consider it does not exist
	return false
}

func updateCaddyConfigForMainCategory(domain, mainCategoryNumber string) error {
	frontendPathsMap := map[string][]string{
		"1":  {"http://localhost:8084"},
		"14": {"http://localhost:8084"},
		"21": {"http://localhost:8084"},
		"23": {"http://localhost:8084"},
		"29": {"http://localhost:8084"},
	}

	frontendPaths, ok := frontendPathsMap[mainCategoryNumber]
	if !ok {
		return fmt.Errorf("invalid main category number: %s", mainCategoryNumber)
	}

	reverseProxies := make([]string, len(frontendPaths))
	for i, path := range frontendPaths {
		reverseProxies[i] = fmt.Sprintf("reverse_proxy %s", path)
	}

	// اضافه کردن import cors به پیکربندی
	newConfig := fmt.Sprintf(`%s {
    %s
    import cors
}
`, domain, strings.Join(reverseProxies, "\n    "))

	return updateCaddyConfig(newConfig)
}

func updateCaddyConfig(newConfig string) error {
	// Open the file for appending, create if it doesn't exist
	file, err := os.OpenFile(caddyfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add a comment to indicate the start of the new configuration
	if _, err := file.WriteString("\n# New WebAR configuration\n"); err != nil {
		return err
	}

	// Write the updated content
	if _, err := file.WriteString(newConfig); err != nil {
		return err
	}

	// Add a comment to indicate the end of the new configuration
	if _, err := file.WriteString("# End of new WebAR configuration\n"); err != nil {
		return err
	}

	// Reload Caddy to apply changes
	cmd := exec.Command("sudo", "systemctl", "reload", "caddy")
	return cmd.Run()
}
