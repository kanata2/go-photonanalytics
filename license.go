package photonanalytics

import "image"

type GetLicenseGraphRequest struct {
	Hash     string
	Template string
	Start    string
	End      string
	Width    int
	Height   int
}

// TODO: implement later
func (c *APIClient) GetLicenseGraph(req *GetLicenseGraphRequest) (image.Image, error) {
	return nil, nil
}

type GetLicenseValueRequest struct {
	Hash     string
	Template string
	Start    string
	End      string
}

// TODO: implement later
func (c *APIClient) GetLicenseValue(req *GetLicenseValueRequest) (float64, error) {
	return 0.0, nil
}

type GetMultipleLicensesValueRequest struct {
	Queries []LicenseBulkDataRequestQuery
}

type LicenseBulkDataRequestQuery struct {
	ID       string
	Hash     string
	Template string
	Start    string
	End      string
}

// TODO: implement later
func (c *APIClient) GetMultipleLicensesValue(req *GetMultipleLicensesValueRequest) (map[string]float64, error) {
	return nil, nil
}
