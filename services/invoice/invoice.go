package invoice

import (
	"path/filepath"

	http_transports "gitag.ir/armogroup/armo/services/reality/services/invoice/http"
	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, amqpClient *amqp.Connection, logger log.Logger, prefix string) service.Invoice {
	invoice := service.MakeInvoice(db, logger)

	httpResource := http_transports.MakeResource(invoice, logger)
	httpResource.RegisterHandlers(r, filepath.Join("/api", prefix))

	// _ = grpc_transports.MakeResource(invoice, logger)
	// grpcResource.RegisterHandlers(r, filepath.Join("/api", prefix))

	// InvoiceQueueName := "invoice.process"
	// agent := amqp_transports.MakeAgent(amqpClient, InvoiceQueueName, logger, invoice)
	// agent.RegisterConsumers()

	return invoice
}
