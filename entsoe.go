//go:generate go run ./tools/extract_types.go

package entsoe

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	periodLayout = "200601021504"

	// requestTimeout bounds a single call to the Transparency Platform. Without
	// it the package used http.DefaultClient, which has no timeout: during the
	// September 2026 maintenance the platform stopped answering rather than
	// returning 503, and every caller blocked until the OS gave up on the TCP
	// connection. A caller sweeping several areas stalls for as long as the
	// platform hangs, so the failure has to be bounded here.
	requestTimeout = 30 * time.Second

	// maxErrorBodyBytes caps how much of an unparseable response is quoted back
	// in an error. The platform answers maintenance with a full HTML page, and
	// embedding it whole put ~300 KB into a single log line.
	maxErrorBodyBytes = 512
)

type EntsoeClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewEntsoeClient(apiKey string) *EntsoeClient {
	c := EntsoeClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: requestTimeout},
	}
	return &c
}

func NewEntsoeClientFromEnv() *EntsoeClient {

	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		logger.Warn().Err(err).Msg("Error loading .env file")
	}

	apiKey := os.Getenv("ENTSOE_API_KEY")
	if apiKey == "" {
		logger.Fatal().Msg("Environment variable ENTSOE_API_KEY with api key not set")
	}

	c := EntsoeClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: requestTimeout},
	}
	return &c
}

// Helper functions

var durations = map[ResolutionType]time.Duration{
	ResolutionQuarter:  15 * time.Minute,
	ResolutionHalfHour: 30 * time.Minute,
	ResolutionHour:     time.Hour,
}

func GetPointTime(start time.Time, position int, resolution ResolutionType) time.Time {
	offset := position - 1

	switch resolution {
	case ResolutionQuarter, ResolutionHalfHour, ResolutionHour:
		return start.Add(time.Duration(offset) * durations[resolution])
	case ResolutionDay:
		return start.AddDate(0, 0, offset)
	case ResolutionWeek:
		return start.AddDate(0, 0, 7*offset)
	case ResolutionYear:
		return start.AddDate(offset, 0, 0)
	}

	return time.Time{}
}

// 4.1. Load domain

// 4.1.1. Actual Total Load [6.1.A]
func (c *EntsoeClient) GetActualTotalLoad(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeSystemTotalLoad))
	params.Add(ParameterProcessType, string(ProcessTypeRealised))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.1.2. Day-Ahead Total Load Forecast [6.1.B]
func (c *EntsoeClient) GetDayAheadTotalLoadForecast(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeSystemTotalLoad))
	params.Add(ParameterProcessType, string(ProcessTypeDayAhead))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.1.3. Week-Ahead Total Load Forecast [6.1.C]
func (c *EntsoeClient) GetWeekAheadTotalLoadForecast(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeSystemTotalLoad))
	params.Add(ParameterProcessType, string(ProcessTypeWeekAhead))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.1.4. Month-Ahead Total Load Forecast [6.1.D]
func (c *EntsoeClient) GetMonthAheadTotalLoadForecast(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeSystemTotalLoad))
	params.Add(ParameterProcessType, string(ProcessTypeMonthAhead))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.1.5. Year-Ahead Total Load Forecast [6.1.E]
func (c *EntsoeClient) GetYearAheadTotalLoadForecast(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeSystemTotalLoad))
	params.Add(ParameterProcessType, string(ProcessTypeYearAhead))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.1.6. Year-Ahead Forecast Margin [8.1]
func (c *EntsoeClient) GetYearAheadForecastMargin(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeLoadForecastMargin))
	params.Add(ParameterProcessType, string(ProcessTypeYearAhead))
	params.Add(ParameterOutBiddingZoneDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.2. Transmission domain

// 4.2.1. Expansion and Dismantling Projects [9.1]
func (c *EntsoeClient) GetExpansionAndDismantlingProjects(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	business *BusinessType,
	docStatus *DocStatus,
) (*TransmissionNetworkMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeInterconnectionNetworkExpansion))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if business != nil {
		params.Add(ParameterBusinessType, string(*business))
	}
	if docStatus != nil {
		params.Add(ParameterDocStatus, string(*docStatus))
	}
	return c.requestTransmissionNetworkMarketDocument(params)
}

// 4.2.2. Forecasted Capacity [11.1.A]
func (c *EntsoeClient) GetForecastedCapacity(
	contractMarketAgreement ContractMarketAgreementType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeEstimatedNetTransferCapacity))
	params.Add(ParameterContractMarketAgreementType, string(contractMarketAgreement))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.3. Offered Capacity [11.1.A]
func (c *EntsoeClient) GetOfferedCapacity(
	auctionType AuctionType,
	contractMarketAgreement ContractMarketAgreementType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	auctionCategory *AuctionCategory,
	classificationSequenceAttributeInstanceComponentPosition *int,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeAgreedCapacity))
	params.Add(ParameterAuctionType, string(auctionType))
	params.Add(ParameterContractMarketAgreementType, string(contractMarketAgreement))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if auctionCategory != nil {
		params.Add(ParameterAuctionCategory, string(*auctionCategory))
	}
	if auctionCategory != nil {
		params.Add(ParameterClassificationSequenceAttributeInstanceComponentPosition, strconv.Itoa(*classificationSequenceAttributeInstanceComponentPosition))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.2.4. Flow-based Parameters [11.1.B]
func (c *EntsoeClient) GetFlowBasedParameters(
	processType ProcessType,
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*CriticalNetworkElementMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeFlowBasedAllocations))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(domain))
	params.Add(ParameterOutDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestCriticalNetworkElementMarketDocument(params)
}

// 4.2.5. Intraday Transfer Limits [11.3]
func (c *EntsoeClient) GetIntradayTransferLimits(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeDcLinkCapacity))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.6. Explicit Allocation Information (Capacity) [12.1.A]
// 4.2.7. Explicit Allocation Information (Revenue only) [12.1.A]
func (c *EntsoeClient) GetExplicitAllocationInformation(
	businessType BusinessType,
	contractMarketAgreementType ContractMarketAgreementType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	auctionCategory *AuctionCategory,
	classificationSequenceAttributeInstanceComponentPosition *int,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeAllocationResultDocument))
	params.Add(ParameterBusinessType, string(businessType))
	params.Add(ParameterContractMarketAgreementType, string(contractMarketAgreementType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if auctionCategory != nil {
		params.Add(ParameterAuctionCategory, string(*auctionCategory))
	}
	if classificationSequenceAttributeInstanceComponentPosition != nil {
		params.Add(ParameterClassificationSequenceAttributeInstanceComponentPosition, strconv.Itoa(*classificationSequenceAttributeInstanceComponentPosition))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.2.8. Total Capacity Nominated [12.1.B]
func (c *EntsoeClient) GetTotalCapacityNominated(
	businessType BusinessType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeCapacityDocument))
	params.Add(ParameterBusinessType, string(businessType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.9. Total Capacity Already Allocated [12.1.C]
func (c *EntsoeClient) GetTotalCapacityAlreadyAllocated(
	businessType BusinessType,
	contractMarketAgreementType ContractMarketAgreementType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	auctionCategory *AuctionCategory,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeCapacityDocument))
	params.Add(ParameterBusinessType, string(businessType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if auctionCategory != nil {
		params.Add(ParameterAuctionCategory, string(*auctionCategory))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.2.10. Day Ahead Prices [12.1.D]
// GetDayAheadPrices returns SDAC (Single Day-Ahead Coupling) prices.
// Resolution was PT60M before 2025-10-01, PT15M after, aligned on Europe/Amsterdam business days.
func (c *EntsoeClient) GetDayAheadPrices(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypePriceDocument))
	params.Add(ParameterContractMarketAgreementType, string(ContractMarketAgreementTypeDaily))
	params.Add(ParameterInDomain, string(domain))
	params.Add(ParameterOutDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.11. Implicit Auction — Net Positions [12.1.E]
// 4.2.12. Implicit Auction — Congestion Income [12.1.E]
func (c *EntsoeClient) GetImplicitAuction(
	businessType BusinessType,
	contractMarketAgreementType ContractMarketAgreementType,
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeAllocationResultDocument))
	params.Add(ParameterBusinessType, string(businessType))
	params.Add(ParameterContractMarketAgreementType, string(contractMarketAgreementType))
	params.Add(ParameterInDomain, string(domain))
	params.Add(ParameterOutDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.13. Total Commercial Schedules [12.1.F]
func (c *EntsoeClient) GetTotalCommercialSchedules(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	contractType *ContractMarketAgreementType,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeFinalisedSchedule))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if contractType != nil {
		params.Add(ParameterContractMarketAgreementType, string(*contractType))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.2.14. Day-ahead Commercial Schedules [12.1.F]
func (c *EntsoeClient) GetDayAheadCommercialSchedules(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	contractType *ContractMarketAgreementType,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeFinalisedSchedule))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if contractType != nil {
		params.Add(ParameterContractMarketAgreementType, string(*contractType))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.2.15. Physical Flows [12.1.G]
func (c *EntsoeClient) GetPhysicalFlows(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeAggregatedEnergyDataReport))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestPublicationMarketDocument(params)
}

// 4.2.16. Capacity Allocated Outside EU [12.1.H]
func (c *EntsoeClient) GetCapacityAllocatedOutsideEu(
	auctionType AuctionType,
	contractMarketAgreementType ContractMarketAgreementType,
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	auctionCategory *AuctionCategory,
	classificationSequenceAttributeInstanceComponentPosition *int,
) (*PublicationMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeNonEuAllocations))
	params.Add(ParameterAuctionType, string(auctionType))
	params.Add(ParameterContractMarketAgreementType, string(contractMarketAgreementType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if auctionCategory != nil {
		params.Add(ParameterAuctionCategory, string(*auctionCategory))
	}
	if classificationSequenceAttributeInstanceComponentPosition != nil {
		params.Add(ParameterClassificationSequenceAttributeInstanceComponentPosition, strconv.Itoa(*classificationSequenceAttributeInstanceComponentPosition))
	}
	return c.requestPublicationMarketDocument(params)
}

// 4.3. Congestion domain

// 4.3.1. Redispatching [13.1.A]
func (c *EntsoeClient) GetRedispatching(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	business *BusinessType,
) (*TransmissionNetworkMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeRedispatchNotice))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if business != nil {
		params.Add(ParameterBusinessType, string(*business))
	}
	return c.requestTransmissionNetworkMarketDocument(params)
}

// 4.3.2. Countertrading [13.1.B]
func (c *EntsoeClient) GetCountertrading(
	inDomain DomainType,
	outDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*TransmissionNetworkMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeCounterTradeNotice))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterOutDomain, string(outDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestTransmissionNetworkMarketDocument(params)
}

// 4.3.3. Costs of Congestion Management [13.1.C]
func (c *EntsoeClient) GetCostsOfCongestionManagement(
	domain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	business *BusinessType,
) (*TransmissionNetworkMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeCongestionCosts))
	params.Add(ParameterInDomain, string(domain))
	params.Add(ParameterOutDomain, string(domain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if business != nil {
		params.Add(ParameterBusinessType, string(*business))
	}
	return c.requestTransmissionNetworkMarketDocument(params)
}

// 4.4. Generation domain

// 4.4.1. Installed Generation Capacity Aggregated [14.1.A]
func (c *EntsoeClient) GetInstalledGenerationCapacityAggregated(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	psrType *PsrType,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeInstalledGenerationPerType))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if psrType != nil {
		params.Add(ParameterPsrType, string(*psrType))
	}
	return c.requestGLMarketDocument(params)
}

// 4.4.2. Installed Generation Capacity per Unit [14.1.B]
func (c *EntsoeClient) GetInstalledGenerationCapacityPerUnit(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	psrType *PsrType,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeGenerationForecast))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if psrType != nil {
		params.Add(ParameterPsrType, string(*psrType))
	}
	return c.requestGLMarketDocument(params)
}

// 4.4.3. Day-ahead Aggregated Generation [14.1.C]
func (c *EntsoeClient) GetDayAheadAggregatedGeneration(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeGenerationForecast))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.4.4. Day-ahead Generation Forecasts for Wind and Solar [14.1.D]
// 4.4.5. Current Generation Forecasts for Wind and Solar [14.1.D]
// 4.4.6. Intraday Generation Forecasts for Wind and Solar [14.1.D]
func (c *EntsoeClient) GetGenerationForecastsForWindAndSolar(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	psrType *PsrType,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeWindAndSolarForecast))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if psrType != nil {
		params.Add(ParameterPsrType, string(*psrType))
	}
	return c.requestGLMarketDocument(params)
}

// 4.4.7. Actual Generation Output per Generation Unit [16.1.A]
func (c *EntsoeClient) GetActualGenerationOutputPerGenerationUnit(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
	psrType *PsrType,
	registeredResource *string,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeActualGeneration))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	if psrType != nil {
		params.Add(ParameterPsrType, string(*psrType))
	}
	if registeredResource != nil {
		params.Add(ParameterRegisteredResource, *registeredResource)
	}
	return c.requestGLMarketDocument(params)
}

// 4.4.8. Aggregated Generation per Type [16.1.B&C]
func (c *EntsoeClient) GetAggregatedGenerationPerType(
	processType ProcessType,
	psrType PsrType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeActualGenerationPerType))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterPsrType, string(psrType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

// 4.4.9. Aggregated Filling Rate of Water Reservoirs and Hydro Storage Plants [16.1.D]
func (c *EntsoeClient) GetAggregatedFillingRateOfWaterReservoirsAndHydroStoragePlants(
	processType ProcessType,
	inDomain DomainType,
	periodStart time.Time,
	periodEnd time.Time,
) (*GLMarketDocument, error) {
	params := url.Values{}
	params.Add(ParameterDocumentType, string(DocumentTypeReservoirFillingInformation))
	params.Add(ParameterProcessType, string(processType))
	params.Add(ParameterInDomain, string(inDomain))
	params.Add(ParameterPeriodStart, periodStart.UTC().Format(periodLayout))
	params.Add(ParameterPeriodEnd, periodEnd.UTC().Format(periodLayout))
	return c.requestGLMarketDocument(params)
}

func (c *EntsoeClient) requestGLMarketDocument(params url.Values) (*GLMarketDocument, error) {
	paramStr := params.Encode()
	data, err := c.sendRequest(paramStr)
	if err != nil {
		return nil, err
	}
	var doc GLMarketDocument
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		return nil, generateParsingError(data)
	}
	return &doc, nil
}

func (c *EntsoeClient) requestTransmissionNetworkMarketDocument(params url.Values) (*TransmissionNetworkMarketDocument, error) {
	paramStr := params.Encode()
	data, err := c.sendRequest(paramStr)
	if err != nil {
		return nil, err
	}

	var doc TransmissionNetworkMarketDocument
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		return nil, generateParsingError(data)
	}
	return &doc, nil
}

func (c *EntsoeClient) requestPublicationMarketDocument(params url.Values) (*PublicationMarketDocument, error) {
	paramStr := params.Encode()
	data, err := c.sendRequest(paramStr)
	if err != nil {
		return nil, err
	}

	var doc PublicationMarketDocument
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		return nil, generateParsingError(data)
	}
	return &doc, nil
}

func (c *EntsoeClient) requestCriticalNetworkElementMarketDocument(params url.Values) (*CriticalNetworkElementMarketDocument, error) {
	paramStr := params.Encode()
	data, err := c.sendRequest(paramStr)
	if err != nil {
		return nil, err
	}

	var doc CriticalNetworkElementMarketDocument
	err = xml.Unmarshal(data, &doc)
	if err != nil {
		return nil, generateParsingError(data)
	}
	return &doc, nil
}

func generateParsingError(data []byte) error {
	d, err := parseAcknowledgementMarketDocument(data)
	if err != nil {
		return err
	}
	return fmt.Errorf("Error requesting data: %s", d.Reason.Text)
}

func parseAcknowledgementMarketDocument(data []byte) (*AcknowledgementMarketDocument, error) {
	var doc AcknowledgementMarketDocument
	err := xml.Unmarshal(data, &doc)
	if err != nil {
		return nil, fmt.Errorf("Error parsing AcknowledgementMarketDocument: %w\n%s", err, truncateForError(data))
	}
	return &doc, nil
}

// redactAPIKey rewrites an error so the API key never reaches a log. The
// message is flattened in the process, which is acceptable: callers log these
// errors, they do not match on wrapped types.
func (c *EntsoeClient) redactAPIKey(err error) error {
	if err == nil || c.apiKey == "" {
		return err
	}
	return fmt.Errorf("%s", strings.ReplaceAll(err.Error(), c.apiKey, "REDACTED"))
}

// truncateForError renders a response body for inclusion in an error message,
// keeping it short enough to stay readable in a log line.
func truncateForError(data []byte) string {
	if len(data) <= maxErrorBodyBytes {
		return string(data)
	}
	return fmt.Sprintf("%s… (%d bytes total)", data[:maxErrorBodyBytes], len(data))
}

func (c *EntsoeClient) sendRequest(paramStr string) ([]byte, error) {
	// A zero-value EntsoeClient built outside the constructors would carry no
	// client at all, so fall back to a bounded one rather than to the unbounded
	// http.DefaultClient.
	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: requestTimeout}
	}

	resp, err := httpClient.Get("https://web-api.tp.entsoe.eu/api?securityToken=" + c.apiKey + "&" + paramStr)
	if err != nil {
		// net/http quotes the full URL in its errors, and the API key travels in
		// the query string, so an unredacted error writes the credential into
		// every caller's logs.
		return nil, c.redactAPIKey(err)
	}
	body := resp.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
