package constants

import "be-b-impact.com/csr/model"

var UserSeed = []model.User{
	{
		Email:    "superadmin@mail.com",
		Username: "superadmin",
		Password: "$2a$10$VQxYnfvUoCx7sGL342w54ex.Z4GFNZO6S/.IJToc22E59hk55OysC",
		Role:     "super",
		Status:   "1",
	},
	{
		Email:    "admin@mail.com",
		Username: "admin1",
		Password: "$2a$10$VQxYnfvUoCx7sGL342w54ex.Z4GFNZO6S/.IJToc22E59hk55OysC",
		Role:     "admin",
		Status:   "1",
	},
}

var TagSeed = []model.Tag{
	{
		Name: "social",
	},
	{
		Name: "volunteer",
	},
	{
		Name: "health",
	},
}

var ProgressSeed = []model.Progress{
	{
		Name:  "Submission Received",
		Label: "received",
	},
	{
		Name:  "Under Review",
		Label: "review",
	},
	{
		Name:  "Approve or Reject",
		Label: "decision",
	},
	{
		Name:  "Proposal Approved",
		Label: "approved",
	},
	{
		Name:  "Proposal Rejected",
		Label: "rejected",
	},
	{
		Name:  "Project In Progress",
		Label: "inProgress",
	},
	{
		Name:  "Project Completed",
		Label: "completed",
	},
}

var CategorySeed = []model.Category{
	{
		Parent:    "0",
		UseFor:    "content",
		Name:      "News",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "content",
		Name:      "Article",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "content",
		Name:      "Report",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "partnership",
		Name:      "Funding",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "partnership",
		Name:      "Event",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "organization",
		Name:      "Non-Profit",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "organization",
		Name:      "Educational Intitution",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "organization",
		Name:      "Healthcare",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "organization",
		Name:      "Others",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Competition",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Exhibition",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Social Actions",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Sport",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Culture",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Seminar",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "event",
		Name:      "Campaign",
		Status:    "1",
		CreatedBy: "super",
	},
	{
		Parent:    "0",
		UseFor:    "general",
		Name:      "General",
		Status:    "1",
		CreatedBy: "super",
	},
}
