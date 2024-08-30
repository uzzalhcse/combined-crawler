package osg

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "osg",
		URL:  "https://osg.icata.net/product-search/catalog/item/",
		Engine: ninjacrawler.Engine{
			IsDynamic: ninjacrawler.Bool(false),
			//BoostCrawling:   true,
			DevCrawlLimit:   15,
			ConcurrentLimit: 4,
			SleepAfter:      50,
			Timeout:         30,
		},
		Preference: ninjacrawler.AppPreference{ExcludeUniqueUrlEntities: []string{constant.ProductDetails}},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}

/*

https://osg.icata.net/product-search/catalog/item/?tool-group=%E3%82%BF%E3%83%83%E3%83%97&abbreviation_method=partial_match&type_for_tap=&screw_thread_symbol=&overall_length[value]=&overall_length[fromto]=&thread_length[value]=&thread_length[fromto]=&shank_diameter[value]=&shank_diameter[fromto]=&work-material_for_tap=&a-brand=&g-list_no=&category_for_tap=&page_no=1
https://osg.icata.net/product-search/catalog/item/?tool-group=%E3%82%BF%E3%83%83%E3%83%97&edp_no=&abbreviation=&abbreviation_method=partial_match&tool-material_for_tap=&type_for_tap=&screw_thread_symbol=&overall_length[value]=&overall_length[fromto]=&thread_length[value]=&thread_length[fromto]=&shank_diameter[value]=&shank_diameter[fromto]=&work-material_for_tap=&a-brand=&g-list_no=&category_for_tap=&page_no=1
https://osg.icata.net/product-search/catalog/item/?tool-group=%E3%83%89%E3%83%AA%E3%83%AB&edp_no=&abbreviation=&abbreviation_method=partial_match&tool-material_for_drill=&drill_diameter=&overall_length[value]=&overall_length[fromto]=&flute_length[value]=&flute_length[fromto]=&shank_diameter[value]=&shank_diameter[fromto]=&work-material_for_drill=&a-brand=&g-list_no=&category_for_drill=&page_no=1#!/search-category-criteria-002/
https://osg.icata.net/product-search/catalog/item/?tool-group=%E3%82%BF%E3%83%83%E3%83%97&edp_no=&abbreviation=&abbreviation_method=partial_match&tool-material_for_tap=&type_for_tap=&screw_thread_symbol=&overall_length[value]=&overall_length[fromto]=&thread_length[value]=&thread_length[fromto]=&shank_diameter[value]=&shank_diameter[fromto]=&work-material_for_tap=&a-brand=&g-list_no=&category_for_tap=&page_no=2
*/
