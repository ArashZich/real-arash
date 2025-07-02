package endpoints

import (
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

type SetupShowroomUrlRequest struct {
	ShowroomUrl        string `json:"showroom_url"`
	MainCategoryNumber int    `json:"mainCategoryNumber"`
}

func (c *SetupShowroomUrlRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	// Add validation logic for the SetupShowroomUrlRequest struct
	return nil
}

func (s *service) SetupShowroomUrl(ctx context.Context, input SetupShowroomUrlRequest) response.ErrorResponse {
	// Remove "https://" prefix from the domain
	showroomUrl := strings.TrimPrefix(input.ShowroomUrl, "https://")

	// Check if the showroomUrl is already registered
	if domainExists(showroomUrl) {
		return response.ErrorBadRequest("دامنه قبلا ثبت شده است")
	}

	// Check if the showroomUrl IP matches the server IP
	ips, err := net.LookupIP(showroomUrl)
	if err != nil {
		log.Printf("Error looking up IP for showroomUrl %s: %v", showroomUrl, err)
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
	log.Printf("Received SetupShowroomUrlRequest: %+v\n", showroomUrl)

	// Implement the logic to update Caddy configuration based on the input
	err = updateCaddyConfigForMainCategoryShowroom(showroomUrl, fmt.Sprint(input.MainCategoryNumber))
	if err != nil {
		s.logger.With(ctx).Error(err)
		return response.ErrorNotFound("خطا در تنظیم دامنه")
	}

	return response.ErrorResponse{}
}

func updateCaddyConfigForMainCategoryShowroom(domain, mainCategoryNumber string) error {
	frontendPathsMap := map[string][]string{
		"1":  {"http://localhost:8086"},
		"14": {"http://localhost:8086"},
		"21": {"http://localhost:8086"},
		"23": {"http://localhost:8086"},
		"29": {"http://localhost:8086"},
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

	return updateCaddyConfigShowroom(newConfig)
}

func updateCaddyConfigShowroom(newConfig string) error {
	// Open the file for appending, create if it doesn't exist
	file, err := os.OpenFile(caddyfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add a comment to indicate the start of the new configuration
	if _, err := file.WriteString("\n# New Showroom configuration\n"); err != nil {
		return err
	}

	// Write the updated content
	if _, err := file.WriteString(newConfig); err != nil {
		return err
	}

	// Add a comment to indicate the end of the new configuration
	if _, err := file.WriteString("# End of new Showroom configuration\n"); err != nil {
		return err
	}

	// Reload Caddy to apply changes
	cmd := exec.Command("sudo", "systemctl", "reload", "caddy")
	return cmd.Run()
}
