package product

type CreateProductRequest struct {
	LevelID     uint    `json:"level_id" binding:"required,gt=0"`
	StepID      uint    `json:"step_id" binding:"required,gt=0"`
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Rating      float64 `json:"rating" binding:"omitempty,gte=0,lte=5"`
	MinQuantity int     `json:"min_quantity" binding:"required,gt=0"`
	MaxQuantity int     `json:"max_quantity" binding:"required,gt=0,gtefield=MinQuantity"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
	Description string  `json:"description" binding:"omitempty"`
}

type UpdateProductRequest struct {
	LevelID     *uint    `json:"level_id" binding:"omitempty,gt=0"`
	StepID      *uint    `json:"step_id" binding:"omitempty,gt=0"`
	Name        *string  `json:"name" binding:"omitempty,min=2,max=255"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Rating      *float64 `json:"rating" binding:"omitempty,gte=0,lte=5"`
	MinQuantity *int     `json:"min_quantity" binding:"omitempty,gt=0"`
	MaxQuantity *int     `json:"max_quantity" binding:"omitempty,gt=0,gtefield=MinQuantity"`
	ImageURL    *string  `json:"image_url" binding:"omitempty,url"`
	Description *string  `json:"description" binding:"omitempty"`
}

type ProductListResponse struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

type LevelResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type StepResponse struct {
	ID         uint   `json:"id"`
	StepNumber int    `json:"step_number"`
	Name       string `json:"name"`
}

type ProductResponse struct {
	ID          string         `json:"id"`
	LevelID     uint           `json:"level_id"`
	Level       *LevelResponse `json:"level,omitempty"`
	StepID      uint           `json:"step_id"`
	Step        *StepResponse  `json:"step,omitempty"`
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	Rating      float64        `json:"rating"`
	MinQuantity int            `json:"min_quantity"`
	MaxQuantity int            `json:"max_quantity"`
	ImageURL    string         `json:"image_url"`
	Description string         `json:"description"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

func ToResponse(p *Product) *ProductResponse {
	resp := &ProductResponse{
		ID:          p.ID,
		LevelID:     p.LevelID,
		StepID:      p.StepID,
		Name:        p.Name,
		Price:       p.Price,
		Rating:      p.Rating,
		MinQuantity: p.MinQuantity,
		MaxQuantity: p.MaxQuantity,
		ImageURL:    p.ImageURL,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if p.Level != nil {
		resp.Level = &LevelResponse{
			ID:   p.Level.ID,
			Name: p.Level.Name,
		}
	}

	if p.Step != nil {
		resp.Step = &StepResponse{
			ID:         p.Step.ID,
			StepNumber: p.Step.StepNumber,
			Name:       p.Step.Name,
		}
	}

	return resp
}

func ToResponseList(products []*Product) []*ProductResponse {
	responses := make([]*ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToResponse(p)
	}
	return responses
}

func ToLevelResponseList(levels []*Level) []*LevelResponse {
	responses := make([]*LevelResponse, len(levels))
	for i, l := range levels {
		responses[i] = &LevelResponse{
			ID:   l.ID,
			Name: l.Name,
		}
	}
	return responses
}

func ToStepResponseList(steps []*Step) []*StepResponse {
	responses := make([]*StepResponse, len(steps))
	for i, s := range steps {
		responses[i] = &StepResponse{
			ID:         s.ID,
			StepNumber: s.StepNumber,
			Name:       s.Name,
		}
	}
	return responses
}
