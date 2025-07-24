package postgresql

import (
	"log"
	"time"

	"github.com/Ablebil/eco-sample/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := seedBadges(db); err != nil {
		return err
	}

	if err := seedChallenges(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedBadges(db *gorm.DB) error {
	log.Println("Seeding badges...")

	badges := []entity.Badge{
		{
			Type:        entity.BadgeEcoWarrior,
			Name:        "Eco Warrior",
			Description: stringPtr("Complete your first challenge and start your eco-friendly journey!"),
			ImageURL:    stringPtr("https://example.com/images/badges/eco-warrior.png"),
			RequiredExp: 50,
		},
		{
			Type:        entity.BadgeGreenHero,
			Name:        "Green Hero",
			Description: stringPtr("You're making a real difference! Keep up the great work."),
			ImageURL:    stringPtr("https://example.com/images/badges/green-hero.png"),
			RequiredExp: 200,
		},
		{
			Type:        entity.BadgeZeroEmission,
			Name:        "Zero Emission Champion",
			Description: stringPtr("Outstanding commitment to reducing carbon footprint."),
			ImageURL:    stringPtr("https://example.com/images/badges/zero-emission.png"),
			RequiredExp: 500,
		},
		{
			Type:        entity.BadgeClimateChampion,
			Name:        "Climate Champion",
			Description: stringPtr("Ultimate eco-warrior! You're a true climate champion."),
			ImageURL:    stringPtr("https://example.com/images/badges/climate-champion.png"),
			RequiredExp: 1000,
		},
	}

	for _, badge := range badges {
		var existingBadge entity.Badge
		err := db.Where("type = ?", badge.Type).First(&existingBadge).Error

		if err == gorm.ErrRecordNotFound {
			id, _ := uuid.NewV7()
			badge.ID = id

			now := time.Now()
			badge.CreatedAt = &now

			if err := db.Create(&badge).Error; err != nil {
				log.Printf("Error creating badge %s: %v", badge.Name, err)
				return err
			}
			log.Printf("Created badge: %s (ID: %s)", badge.Name, badge.ID)
		} else if err != nil {
			log.Printf("Error checking badge %s: %v", badge.Name, err)
			return err
		} else {
			log.Printf("Badge %s already exists, skipping", badge.Name)
		}
	}

	return nil
}

func seedChallenges(db *gorm.DB) error {
	log.Println("Seeding challenges...")

	challenges := []entity.Challenge{
		{
			Title:       "Meatless Monday",
			Description: stringPtr("Go vegetarian for a full day. Skip meat and try delicious plant-based alternatives!"),
			ExpReward:   25,
			IsActive:    true,
		},
		{
			Title:       "Bike to Work",
			Description: stringPtr("Cycle to work instead of using motorized transport. Great for health and environment!"),
			ExpReward:   30,
			IsActive:    true,
		},
		{
			Title:       "Zero Plastic Day",
			Description: stringPtr("Avoid single-use plastics for an entire day. Bring your own bags and containers!"),
			ExpReward:   35,
			IsActive:    true,
		},
		{
			Title:       "Energy Saver",
			Description: stringPtr("Reduce electricity usage by 20% for a day. Unplug devices and use natural light!"),
			ExpReward:   20,
			IsActive:    true,
		},
		{
			Title:       "Water Conservation",
			Description: stringPtr("Implement water-saving techniques for a week. Take shorter showers and fix leaks!"),
			ExpReward:   40,
			IsActive:    true,
		},
		{
			Title:       "Public Transport Champion",
			Description: stringPtr("Use public transportation for all your trips in a day instead of private vehicles."),
			ExpReward:   25,
			IsActive:    true,
		},
		{
			Title:       "Digital Minimalist",
			Description: stringPtr("Reduce screen time and digital consumption for a day. Enjoy offline activities!"),
			ExpReward:   15,
			IsActive:    true,
		},
		{
			Title:       "Local Food Hero",
			Description: stringPtr("Buy only locally sourced food for a week. Support local farmers and reduce transport emissions!"),
			ExpReward:   45,
			IsActive:    true,
		},
		{
			Title:       "Reusable Bottle Week",
			Description: stringPtr("Use only reusable water bottles for a full week. Help reduce plastic waste!"),
			ExpReward:   30,
			IsActive:    true,
		},
		{
			Title:       "Paperless Day",
			Description: stringPtr("Go completely paperless for a day. Use digital alternatives for all documents!"),
			ExpReward:   20,
			IsActive:    true,
		},
	}

	for _, challenge := range challenges {
		var existingChallenge entity.Challenge
		err := db.Where("title = ?", challenge.Title).First(&existingChallenge).Error

		if err == gorm.ErrRecordNotFound {
			id, _ := uuid.NewV7()
			challenge.ID = id

			now := time.Now()
			challenge.CreatedAt = &now
			challenge.UpdatedAt = &now

			if err := db.Create(&challenge).Error; err != nil {
				log.Printf("Error creating challenge %s: %v", challenge.Title, err)
				return err
			}
			log.Printf("Created challenge: %s (ID: %s)", challenge.Title, challenge.ID)
		} else if err != nil {
			log.Printf("Error checking challenge %s: %v", challenge.Title, err)
			return err
		} else {
			log.Printf("Challenge %s already exists, skipping", challenge.Title)
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
