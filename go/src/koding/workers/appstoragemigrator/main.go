package main

import (
	"flag"

	"koding/db/models"
	helper "koding/db/mongodb/modelhelper"
	"koding/helpers"
	"koding/tools/config"
	"koding/tools/logger"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var log = logger.New("appstorage migrator")

var (
	conf           *config.Config
	flagDebug      = flag.Bool("d", false, "Debug mode")
	flagProfile    = flag.String("c", "dev", "Configuration profile from file")
	flagSkip       = flag.Int("s", 0, "Configuration profile from file")
	flagLimit      = flag.Int("l", 1000, "Configuration profile from file")
	collectionName = "jAccounts"
)

func initialize() {
	flag.Parse()
	if *flagProfile == "" {
		log.Fatal("Please specify profile via -c. Aborting.")
	}

	conf = config.MustConfig(*flagProfile)
	helper.Initialize(conf.Mongo)
}

func main() {
	initialize()
	log.Info("worker started")

	var resultDataType models.Account

	iterOptions := helpers.NewIterOptions()
	iterOptions.CollectionName = collectionName
	iterOptions.F = processDocuments
	iterOptions.Filter = createFilter()
	iterOptions.Result = &resultDataType
	iterOptions.Limit = *flagLimit
	iterOptions.Skip = *flagSkip
	iterOptions.Log = log

	log.SetLevel(logger.INFO)
	if *flagDebug {
		log.SetLevel(logger.DEBUG)
	}

	err := helpers.Iter(helper.Mongo, iterOptions)
	if err != nil {
		log.Fatal("Error while iter: %v", err)
	}
	log.Info("worker finished")
}

func createFilter() helper.Selector {
	return helper.Selector{
	// cihangir
	// "_id": bson.ObjectIdHex("50c4a3fe6b33139354000159"),
	// sinan
	// "_id": bson.ObjectIdHex("4f14fa4d519ab4c62e000052"),
	}
}

func processDocuments(doc interface{}) error {
	result := *(doc.(*models.Account))
	rels, err := helper.GetAllRelationships(helper.Selector{
		"targetName": "JAppStorage",
		"sourceName": "JAccount",
		"sourceId":   result.Id,
		"as":         "appStorage",
	})
	if err != nil {
		return err
	}

	// if target user doesnt have any relationship, no need to continue
	if len(rels) == 0 {
		log.Info("couldnt find any relationships")
		return nil
	}

	// aggragate all relationship targetIds, because we are going to fetch target JAppStorages
	ids := make([]bson.ObjectId, 0)
	for _, rel := range rels {
		if rel.TargetId.Valid() {
			ids = append(ids, rel.TargetId)
		}
	}

	// if result doesnt have any valid id, skip
	if len(ids) == 0 {
		log.Info("relationship ids are not valid")
		return nil
	}

	// fetch old JAppStorages
	aps, err := helper.GetAppStoragesByIds(ids...)
	if err != nil {
		return err
	}

	if len(aps) == 0 {
		log.Info("couldnt find any JAppStorages for %s", result.Profile.Nickname)
		return nil
	}

	cas := &models.CombinedAppStorage{
		Id:        bson.NewObjectId(),
		AccountId: result.Id,
		Bucket:    make(map[string]map[string]map[string]interface{}),
	}

	for _, ap := range aps {
		if _, ok := cas.Bucket[ap.AppID]; !ok {
			cas.Bucket[ap.AppID] = make(map[string]map[string]interface{})
		}

		if _, ok := cas.Bucket[ap.AppID]["data"]; !ok {
			cas.Bucket[ap.AppID]["data"] = make(map[string]interface{})
		}

		if len(cas.Bucket[ap.AppID]["data"]) <= len(ap.Bucket) {
			cas.Bucket[ap.AppID]["data"] = ap.Bucket
		} else {
			log.Info("skipping already set apid: %+v", ap.AppID)
		}
	}

	cass, err := helper.GetCombinedAppStorageByAccountId(result.Id)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}

	if err == mgo.ErrNotFound {
		log.Info("creating new CombinedAppStorage for %s", result.Profile.Nickname)
		cas = cleanup(cas)

		if err := helper.CreateCombinedAppStorage(cas); err != nil {
			return err
		}

	} else {
		log.Info("updating existing CombinedAppStorage for %s", result.Profile.Nickname)

		for key, data := range cas.Bucket {
			if _, ok := cass.Bucket[key]; !ok {
				cass.Bucket[key] = data
				continue
			}

			for k, datum := range data {
				if k == "data" {
					continue
				}

				if _, ok := cass.Bucket[key]["data"][k]; !ok {
					cass.Bucket[key]["data"][k] = datum
				}
			}
		}

		cass = cleanup(cass)

		if err := helper.UpdateCombinedAppStorage(cass); err != nil {
			return err
		}
	}

	return nil
}

// cleanup removed blacklisted app storage items, they are mostly not in use
// anymore, list is collected with @sinan
func cleanup(c *models.CombinedAppStorage) *models.CombinedAppStorage {
	if len(c.Bucket) == 0 {
		return c
	}

	for _, black := range blacklist {
		delete(c.Bucket, black)
	}

	return c
}

var blacklist = []string{
	"Activity",
	"Members",
	"WelcomeModal",
	"Wordpress",
	"github-dashboard",
	"laravel-installer",
	"OnboardingStatus",
	"About",
	"AceTabHistory",
	"Applications",
	"Apps",
	"Brackets",
	"Brackets",
	"Bugs",
	"Dashboard",
	"DefaultAppConfig",
	"DevTools",
	"Dock",
	"EnvironmentsScene",
	"Gameoflife",
	"Hartl",
	"Helper",
	"Home",
	"Installer",
	"Groups",
	"IntroductionTooltipStatus",
	"Julia",
	"Kodepad",
	"KodingApps",
	"Login",
	"KodingBook",
	"MainApp",
	"NewKoding",
	"PhoneGap",
	"Numbers",
	"Rubyonrailsinstaller",
	"Teamwork",
	"Topics",
	"Umlgenerator",
}
