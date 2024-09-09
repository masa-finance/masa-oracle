package twitterscraper_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type Response struct {
	Data struct {
		User struct {
			Result struct {
				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Entries []struct {
								Content struct {
									ItemContent struct {
										UserResults struct {
											Result struct {
												Legacy Legacy `json:"legacy"`
											} `json:"result"`
										} `json:"user_results"`
									} `json:"itemContent"`
								} `json:"content"`
							} `json:"entries"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type Legacy struct {
	ScreenName      string `json:"screen_name"`
	FollowersCount  int    `json:"followers_count"`
	FriendsCount    int    `json:"friends_count"`
	ListedCount     int    `json:"listed_count"`
	CreatedAt       string `json:"created_at"`
	FavouritesCount int    `json:"favourites_count"`
	StatusesCount   int    `json:"statuses_count"`
	MediaCount      int    `json:"media_count"`
	ProfileImageUrl string `json:"profile_image_url_https"`
	Description     string `json:"description"`
	Location        string `json:"location"`
	Url             string `json:"url"`
	Protected       bool   `json:"protected"`
	Verified        bool   `json:"verified"`
}

func parseLegacyInfo(jsonString string) ([]Legacy, error) {
	var response Response
	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		return nil, err
	}

	var legacies []Legacy
	for _, instruction := range response.Data.User.Result.Timeline.Timeline.Instructions {
		for _, entry := range instruction.Entries {
			legacies = append(legacies, entry.Content.ItemContent.UserResults.Result.Legacy)
			log.Printf("Added legacy for user: %s\n", entry.Content.ItemContent.UserResults.Result.Legacy.ScreenName) // Log the screen name of the user being added
		}
	}

	return legacies, nil
}

func TestLegacyInfo(t *testing.T) {
	jsonString := `{
		"data": {
			"user": {
				"result": {
					"__typename": "User",
					"timeline": {
						"timeline": {
							"instructions": [
								{
									"type": "TimelineClearCache"
								},
								{
									"type": "TimelineTerminateTimeline",
									"direction": "Top"
								},
								{
									"type": "TimelineTerminateTimeline",
									"direction": "Bottom"
								},
								{
									"type": "TimelineAddEntries",
									"entries": [
										{
											"entryId": "user-16331756",
											"sortIndex": "1790198958568505344",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNjMzMTc1Ng==",
															"rest_id": "16331756",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Square",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Wed Sep 17 16:23:47 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "South Florida's #1 News Station! Your 24/7 source for breaking news, @7Weather & @7SportsXtra powered by our digital team. Breaking news? newsdesk@wsvn.com",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "wsvn.com",
																				"expanded_url": "http://www.wsvn.com",
																				"url": "https://t.co/fx2A1LcmVr",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 4919,
																"followers_count": 479463,
																"friends_count": 867,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2473,
																"location": "Miami / Fort Lauderdale",
																"media_count": 48898,
																"name": "WSVN 7 News",
																"normal_followers_count": 479463,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/16331756/1474555064",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/629323138011123712/vWW4QIUC_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "wsvn",
																"statuses_count": 206704,
																"translator_type": "none",
																"url": "https://t.co/fx2A1LcmVr",
																"verified": false,
																"verified_type": "Business",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1498692248510291970",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 580,
																		"name": "Media & News Company",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-2881087294",
											"sortIndex": "1790198958568505343",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoyODgxMDg3Mjk0",
															"rest_id": "2881087294",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Mon Nov 17 14:19:37 +0000 2014",
																"default_profile": true,
																"default_profile_image": false,
																"description": "First Madam Mayor of Miami-Dade County. Water Warrior and constant fighter for our community, children, and future.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "miamidade.gov/mayor",
																				"expanded_url": "http://www.miamidade.gov/mayor",
																				"url": "https://t.co/aUd5jNbrdJ",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 34557,
																"followers_count": 65756,
																"friends_count": 7116,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 426,
																"location": "Miami-Dade County",
																"media_count": 4705,
																"name": "Daniella Levine Cava",
																"normal_followers_count": 65756,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/2881087294/1707879348",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1659168094625976321/RAsBTKKk_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "MayorDaniella",
																"statuses_count": 17989,
																"translator_type": "none",
																"url": "https://t.co/aUd5jNbrdJ",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1659003107345301504",
																"professional_type": "Business",
																"category": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-253563556",
											"sortIndex": "1790198958568505342",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoyNTM1NjM1NTY=",
															"rest_id": "253563556",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": true,
															"profile_image_shape": "Square",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Thu Feb 17 14:10:20 +0000 2011",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Official Twitter page of the MDPD. EMERGENCIES: DIAL 9-1-1. Non-Emergency: 305-476-5423. Site not monitored 24/7. Terms of Use: https://t.co/eZd89FRI96",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "linktr.ee/MiamiDadePD",
																				"expanded_url": "https://linktr.ee/MiamiDadePD",
																				"url": "https://t.co/eZd89FRI96",
																				"indices": [
																					128,
																					151
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "miamidade.gov/police",
																				"expanded_url": "http://www.miamidade.gov/police",
																				"url": "https://t.co/Feb2V214Ej",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 11165,
																"followers_count": 98246,
																"friends_count": 1014,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 970,
																"location": "Miami-Dade County, Florida",
																"media_count": 17093,
																"name": "Miami-Dade Police",
																"normal_followers_count": 98246,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/253563556/1646248660",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1587435171917217792/1UuHry2U_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "MiamiDadePD",
																"statuses_count": 24875,
																"translator_type": "none",
																"url": "https://t.co/Feb2V214Ej",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-18728203",
											"sortIndex": "1790198958568505341",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxODcyODIwMw==",
															"rest_id": "18728203",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Square",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Wed Jan 07 16:58:28 +0000 2009",
																"default_profile": false,
																"default_profile_image": false,
																"description": "South Florida News, Weather, Entertainment, Sports from WPLG Local 10 Follow us on Instagram: https://t.co/IIfr4MFhzm",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "instagram.com/local10news",
																				"expanded_url": "http://instagram.com/local10news",
																				"url": "https://t.co/IIfr4MFhzm",
																				"indices": [
																					94,
																					117
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "local10.com",
																				"expanded_url": "http://www.local10.com",
																				"url": "https://t.co/fMkL5X3ifE",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3588,
																"followers_count": 234441,
																"friends_count": 1226,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2071,
																"location": "Miami, FL",
																"media_count": 60730,
																"name": "WPLG Local 10 News",
																"normal_followers_count": 234441,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/18728203/1623255556",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1675165975065374720/NmAQ7CZg_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "WPLGLocal10",
																"statuses_count": 189767,
																"translator_type": "none",
																"url": "https://t.co/fMkL5X3ifE",
																"verified": false,
																"verified_type": "Business",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1547665018673451010",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 579,
																		"name": "Media & News",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1704870318",
											"sortIndex": "1790198958568505340",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzA0ODcwMzE4",
															"rest_id": "1704870318",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Tue Aug 27 15:01:31 +0000 2013",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Miami Police Dept. |We are on FB, Ig, Nextdoor, YouTube, and TikTok | Account NOT monitored 24/7 | For Emergencies call 911 | Terms of Use: https://t.co/1O1bfjjV5B",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "goo.gl/eObLVf",
																				"expanded_url": "http://goo.gl/eObLVf",
																				"url": "https://t.co/1O1bfjjV5B",
																				"indices": [
																					140,
																					163
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "miami-police.org/index.asp",
																				"expanded_url": "http://www.miami-police.org/index.asp",
																				"url": "https://t.co/BsG9ue2rrd",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 8972,
																"followers_count": 68297,
																"friends_count": 1431,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 696,
																"location": "Miami, Florida",
																"media_count": 9216,
																"name": "Miami PD",
																"normal_followers_count": 68297,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1704870318/1685294665",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1743442727553937408/1yIiXebt_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "MiamiPD",
																"statuses_count": 16498,
																"translator_type": "none",
																"url": "https://t.co/BsG9ue2rrd",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1634219685792174080",
																"professional_type": "Creator",
																"category": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-617602505",
											"sortIndex": "1790198958568505339",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjo2MTc2MDI1MDU=",
															"rest_id": "617602505",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon Jun 25 00:19:02 +0000 2012",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Proud to serve as @MiamiMayor. Former President of @USMayors.",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 28238,
																"followers_count": 146754,
																"friends_count": 305,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 1042,
																"location": "Miami, FL",
																"media_count": 2678,
																"name": "Mayor Francis Suarez",
																"normal_followers_count": 146754,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1696972096927047680/QxXgUhV__normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "FrancisSuarez",
																"statuses_count": 18092,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1612869099444441088",
																"professional_type": "Creator",
																"category": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-189669352",
											"sortIndex": "1790198958568505338",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxODk2NjkzNTI=",
															"rest_id": "189669352",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": true,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Sat Sep 11 22:23:52 +0000 2010",
																"default_profile": true,
																"default_profile_image": false,
																"description": "\uD83C\uDF34 The Good | The Bad | The Funny \uD83C\uDF34    Citizens Journalism \uD83D\uDDDE️",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "onlyindade.com",
																				"expanded_url": "http://onlyindade.com",
																				"url": "https://t.co/mMddxkvTmX",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1524,
																"followers_count": 102676,
																"friends_count": 567,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 155,
																"location": "Miami Dade County",
																"media_count": 2023,
																"name": "ONLY in DADE",
																"normal_followers_count": 102676,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/189669352/1635442343",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1154710203495002113/4BNNSmr7_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "ONLYinDADE",
																"statuses_count": 7186,
																"translator_type": "none",
																"url": "https://t.co/mMddxkvTmX",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1527694179115364352",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 580,
																		"name": "Media & News Company",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-17931773",
											"sortIndex": "1790198958568505337",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkzMTc3Mw==",
															"rest_id": "17931773",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Square",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Sat Dec 06 23:55:42 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Official account of the City of Miami Beach. Download our Miami Beach Gov app to report any issue.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "miamibeachfl.gov",
																				"expanded_url": "http://www.miamibeachfl.gov",
																				"url": "https://t.co/Pdyb5Qu6ij",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 494,
																"followers_count": 124226,
																"friends_count": 470,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 1179,
																"location": "Fun & Sun Capital of the World",
																"media_count": 17280,
																"name": "City of Miami Beach",
																"normal_followers_count": 124226,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/17931773/1698846656",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1719713635444920320/QHmsnezr_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "MiamiBeachNews",
																"statuses_count": 45268,
																"translator_type": "none",
																"url": "https://t.co/Pdyb5Qu6ij",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-590060558",
											"sortIndex": "1790198958568505336",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjo1OTAwNjA1NTg=",
															"rest_id": "590060558",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Fri May 25 15:03:01 +0000 2012",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Official Twitter account for the National Weather Service Miami-South Florida. Details: https://t.co/fHaf0Ly6pT",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "weather.gov/twitter",
																				"expanded_url": "http://weather.gov/twitter",
																				"url": "https://t.co/fHaf0Ly6pT",
																				"indices": [
																					88,
																					111
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "weather.gov/miami",
																				"expanded_url": "http://www.weather.gov/miami",
																				"url": "https://t.co/ByLOeNUZCN",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 9274,
																"followers_count": 94164,
																"friends_count": 627,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2143,
																"location": "Miami, Florida",
																"media_count": 28287,
																"name": "NWS Miami",
																"normal_followers_count": 94164,
																"pinned_tweet_ids_str": [
																	"1789980155152040423"
																],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/590060558/1695570055",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/876059691616419840/p_BBsv8h_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "NWSMiami",
																"statuses_count": 42307,
																"translator_type": "regular",
																"url": "https://t.co/ByLOeNUZCN",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-16334281",
											"sortIndex": "1790198958568505335",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNjMzNDI4MQ==",
															"rest_id": "16334281",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Wed Sep 17 18:48:22 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "News and information from CBS News Miami.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "CBSNewsMiami.com",
																				"expanded_url": "http://www.CBSNewsMiami.com",
																				"url": "https://t.co/n2vIzLRF15",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 5340,
																"followers_count": 130736,
																"friends_count": 1514,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2150,
																"location": "Miami",
																"media_count": 41983,
																"name": "CBS News Miami",
																"normal_followers_count": 130736,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/16334281/1684766043",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1641116924649127939/e8XCuDzW_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "CBSMiami",
																"statuses_count": 256355,
																"translator_type": "none",
																"url": "https://t.co/n2vIzLRF15",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-15727981",
											"sortIndex": "1790198958568505334",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNTcyNzk4MQ==",
															"rest_id": "15727981",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon Aug 04 21:14:18 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "South Florida's trusted source for Breaking News, Weather and Traffic. Watch us LIVE, 24/7 the #NBC6 app.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "NBC6.com",
																				"expanded_url": "http://www.NBC6.com",
																				"url": "https://t.co/dzHVOTUMBu",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2957,
																"followers_count": 329570,
																"friends_count": 995,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2603,
																"location": "South Florida",
																"media_count": 28583,
																"name": "NBC 6 South Florida",
																"normal_followers_count": 329570,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/15727981/1581633766",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1062821211699376128/HQKreXmc_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "nbc6",
																"statuses_count": 245434,
																"translator_type": "none",
																"url": "https://t.co/dzHVOTUMBu",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1639003450993147908",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 580,
																		"name": "Media & News Company",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-55419844",
											"sortIndex": "1790198958568505333",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjo1NTQxOTg0NA==",
															"rest_id": "55419844",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Fri Jul 10 00:54:19 +0000 2009",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Certified Consulting Meteorologist, ClimaData | Hurricane Specialist @nbc6 | Columnist @bulletinatomic | BSc & Trustee @Cornell | MSc @JohnsHopkins",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "climadata.com",
																				"expanded_url": "http://www.climadata.com",
																				"url": "https://t.co/CEX9Ur0GOr",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 69488,
																"followers_count": 129841,
																"friends_count": 2541,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 1370,
																"location": "Miami Florida USA",
																"media_count": 15153,
																"name": "John Morales",
																"normal_followers_count": 129841,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/55419844/1609945926",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1535619808154337283/JKIyMgPg_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "JohnMoralesTV",
																"statuses_count": 83100,
																"translator_type": "none",
																"url": "https://t.co/CEX9Ur0GOr",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1489287965511663619",
																"professional_type": "Creator",
																"category": [
																	{
																		"id": 933,
																		"name": "Media Personality",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1257764914539659270",
											"sortIndex": "1790198958568505332",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxMjU3NzY0OTE0NTM5NjU5Mjcw",
															"rest_id": "1257764914539659270",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": true,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 05 20:12:23 +0000 2020",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Superintendent @MDCPS 'Every child needs a champion who understands the power of connection and insists they become the best they can be.' R. Pierson",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 16282,
																"followers_count": 10752,
																"friends_count": 728,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 16,
																"location": "",
																"media_count": 1996,
																"name": "Jose L. Dotres, Ed.D.",
																"normal_followers_count": 10752,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1257764914539659270/1700234785",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1504618978924236803/Mnx83D-D_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "SuptDotres",
																"statuses_count": 3614,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1517904564585963520",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 144,
																		"name": "Education",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-12699932",
											"sortIndex": "1790198958568505331",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxMjY5OTkzMg==",
															"rest_id": "12699932",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Square",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Fri Jan 25 22:16:56 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "The official Sun Sentinel Twitter account. Covering Broward, Palm Beach and Miami-Dade. Pulitzer Prize Gold Medal For Public Service 2013, 2019",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "sun-sentinel.com",
																				"expanded_url": "http://www.sun-sentinel.com/",
																				"url": "https://t.co/B1s4KIMA9g",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1971,
																"followers_count": 299523,
																"friends_count": 11552,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2891,
																"location": "South Florida",
																"media_count": 151758,
																"name": "South Florida Sun Sentinel",
																"normal_followers_count": 299523,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/12699932/1712055906",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1073311891303354368/zaX_Xf1G_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "SunSentinel",
																"statuses_count": 234818,
																"translator_type": "regular",
																"url": "https://t.co/B1s4KIMA9g",
																"verified": false,
																"verified_type": "Business",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1488944659820945412",
																"professional_type": "Creator",
																"category": [
																	{
																		"id": 580,
																		"name": "Media & News Company",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-15616612",
											"sortIndex": "1790198958568505330",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNTYxNjYxMg==",
															"rest_id": "15616612",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": true,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Sun Jul 27 02:24:34 +0000 2008",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Miami / #MiamiTech Ambassador & \uD83C\uDFD7️ Opinions & Tweets are my own -\n\uD83D\uDE80 / Advisor at @officelogicmia\nSuper Connector \uD83D\uDC4B",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "office-logic.co",
																				"expanded_url": "https://office-logic.co",
																				"url": "https://t.co/3FEoSOlvP1",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 9468,
																"followers_count": 15090,
																"friends_count": 3213,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 205,
																"location": "Miami, FL",
																"media_count": 21063,
																"name": "Ryan RC Rea",
																"normal_followers_count": 15090,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/15616612/1703717574",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1666169070146797569/0P0NtO-m_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "volvoshine",
																"statuses_count": 54905,
																"translator_type": "none",
																"url": "https://t.co/3FEoSOlvP1",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1460738606050549766",
																"professional_type": "Creator",
																"category": [
																	{
																		"id": 934,
																		"name": "Social Media Influencer",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {
																"is_enabled": true
															}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-156325144",
											"sortIndex": "1790198958568505329",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNTYzMjUxNDQ=",
															"rest_id": "156325144",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Wed Jun 16 16:52:18 +0000 2010",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Chief of Police - Doral Police Dept @DoralPolice @Cityofdoral - former Chief of Police @MDSPD @MDCPS - views are my own/retweets are not endorsements",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 12645,
																"followers_count": 4774,
																"friends_count": 740,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 27,
																"location": "Miami, FL",
																"media_count": 1117,
																"name": "Edwin Lopez",
																"normal_followers_count": 4774,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/156325144/1673068766",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1621330528665325571/26RXZ6zn_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "ChiefDoralPD",
																"statuses_count": 2654,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-135265376",
											"sortIndex": "1790198958568505328",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxMzUyNjUzNzY=",
															"rest_id": "135265376",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue Apr 20 20:28:36 +0000 2010",
																"default_profile": false,
																"default_profile_image": false,
																"description": "The official Twitter of the Miami Beach Police Department. Dial 911 for emergencies or 305.673.7900 for non-emergencies. #YourMBPD",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "yourMBPD.com",
																				"expanded_url": "http://yourMBPD.com",
																				"url": "https://t.co/ZkAwHyQHZj",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 13646,
																"followers_count": 53264,
																"friends_count": 422,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 756,
																"location": " Miami Beach, FL. 33139",
																"media_count": 6008,
																"name": "Miami Beach Police",
																"normal_followers_count": 53264,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/135265376/1706134356",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1750280418421600257/Lo5oN7YC_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "MiamiBeachPD",
																"statuses_count": 19760,
																"translator_type": "none",
																"url": "https://t.co/ZkAwHyQHZj",
																"verified": false,
																"verified_type": "Government",
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-352557959",
											"sortIndex": "1790198958568505327",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjozNTI1NTc5NTk=",
															"rest_id": "352557959",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Wed Aug 10 19:55:03 +0000 2011",
																"default_profile": false,
																"default_profile_image": false,
																"description": "@LASchools Superintendent, National/Urban Superintendent of the Year, McGraw prize winner. Education turned a once homeless kid into me.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "lausd.org/superintendent",
																				"expanded_url": "https://lausd.org/superintendent",
																				"url": "https://t.co/lcf8lS1c6i",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 51945,
																"followers_count": 81698,
																"friends_count": 1208,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 480,
																"location": "Los Angeles, CA",
																"media_count": 7908,
																"name": "Alberto M. Carvalho",
																"normal_followers_count": 81698,
																"pinned_tweet_ids_str": [
																	"1790042092200194462"
																],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/352557959/1703617328",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1708847822194589696/fx_u1nTj_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "LAUSDSup",
																"statuses_count": 30729,
																"translator_type": "none",
																"url": "https://t.co/lcf8lS1c6i",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1372212427",
											"sortIndex": "1790198958568505326",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxMzcyMjEyNDI3",
															"rest_id": "1372212427",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Mon Apr 22 13:56:59 +0000 2013",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Welcome to Miami International Airport (MIA). Monitored M-F, 9am-6pm. After hours \uD83D\uDCDE: (305) 876-7000. \uD83D\uDECD️: @ShopsatMIA",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "iflymia.com",
																				"expanded_url": "http://www.iflymia.com",
																				"url": "https://t.co/kBKbIPPrhQ",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 11233,
																"followers_count": 71351,
																"friends_count": 2261,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 966,
																"location": "Miami, Florida",
																"media_count": 7316,
																"name": "Miami Int'l Airport",
																"normal_followers_count": 71351,
																"pinned_tweet_ids_str": [
																	"1785712667878994065"
																],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1372212427/1713304092",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1742209897125548032/FVN-S4gO_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "iflymia",
																"statuses_count": 33418,
																"translator_type": "none",
																"url": "https://t.co/kBKbIPPrhQ",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1461221109408407553",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 912,
																		"name": "Airport",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-23852881",
											"sortIndex": "1790198958568505325",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoyMzg1Mjg4MQ==",
															"rest_id": "23852881",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Wed Mar 11 23:31:41 +0000 2009",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Doug Hanks covers Miami-Dade County for The Miami Herald. dhanks@miamiherald.com",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "miamiherald.com",
																				"expanded_url": "http://miamiherald.com",
																				"url": "https://t.co/JtOgAeg68X",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 8520,
																"followers_count": 22695,
																"friends_count": 2898,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 473,
																"location": "Miami",
																"media_count": 4967,
																"name": "Doug Hanks",
																"normal_followers_count": 22695,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/23852881/1411047586",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1773455207810109440/Wh_PNo9p_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "doug_hanks",
																"statuses_count": 46265,
																"translator_type": "none",
																"url": "https://t.co/JtOgAeg68X",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790196812587954177",
											"sortIndex": "1790198958568505324",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTk2ODEyNTg3OTU0MTc3",
															"rest_id": "1790196812587954177",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 14 01:46:13 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 0,
																"friends_count": 34,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Jessica Hernández",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790197347261018112/Q0oFG24Z_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "hernandez34720",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790192783120089088",
											"sortIndex": "1790198958568505323",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTkyNzgzMTIwMDg5MDg4",
															"rest_id": "1790192783120089088",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 14 01:29:54 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "You can pay for school but you can't buy class. Gooner through and through.",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 5,
																"followers_count": 0,
																"friends_count": 11,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "Jamesbury, North Dakota",
																"media_count": 0,
																"name": "narasimha",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": true,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790192783120089088/1715650225",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790192895774896128/5Rsu4rvk_normal.jpg",
																"profile_interstitial_type": "offensive_profile_content",
																"screen_name": "herring_jo66714",
																"statuses_count": 2,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790191901271891968",
											"sortIndex": "1790198958568505322",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTkxOTAxMjcxODkxOTY4",
															"rest_id": "1790191901271891968",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 14 01:26:24 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "We make weird games. Follow us to keep up with updates to #MyDadsTower, and any other projects we decide to take on.",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 5,
																"followers_count": 0,
																"friends_count": 11,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "North Gailport, West Virginia",
																"media_count": 0,
																"name": "dadhalla",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790191901271891968/1715650062",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790192215479713792/s14Jf1Dd_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "DanielleWh75538",
																"statuses_count": 2,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790191906489511936",
											"sortIndex": "1790198958568505321",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTkxOTA2NDg5NTExOTM2",
															"rest_id": "1790191906489511936",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 14 01:26:26 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "↗️ Engineering geologist at Simpson-Clay",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 5,
																"followers_count": 0,
																"friends_count": 10,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "Bryanland, Texas",
																"media_count": 0,
																"name": "rainee",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790191906489511936/1715650029",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790192074786017280/dtNjwSp__normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "SmithDonna38928",
																"statuses_count": 2,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1789841120995762177",
											"sortIndex": "1790198958568505320",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzg5ODQxMTIwOTk1NzYyMTc3",
															"rest_id": "1789841120995762177",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 02:13:54 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "“fuck around and find out”",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 3,
																"friends_count": 77,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "Miami, FL",
																"media_count": 1,
																"name": "John doe",
																"normal_followers_count": 3,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1789841120995762177/1715651180",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790197340390752256/Djoipg1z_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "thehimJohndoe",
																"statuses_count": 1,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1052349669358137345",
											"sortIndex": "1790198958568505319",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxMDUyMzQ5NjY5MzU4MTM3MzQ1",
															"rest_id": "1052349669358137345",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": true,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Wed Oct 17 00:04:47 +0000 2018",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Avid sports fan (usually). Falcons, Dawgs, Hawks, Man United, Real Madrid.",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "youtube.com/channel/UC6p0r…",
																				"expanded_url": "https://www.youtube.com/channel/UC6p0r5_AzqhRe15tSx43pNA",
																				"url": "https://t.co/LFTKRaUJej",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 12475,
																"followers_count": 459,
																"friends_count": 1473,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 1,
																"location": "Southeast, United States",
																"media_count": 562,
																"name": "Ultra",
																"normal_followers_count": 459,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1052349669358137345/1673397155",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1787874768898715649/x4xw8kks_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "Ultra_BLV",
																"statuses_count": 6557,
																"translator_type": "none",
																"url": "https://t.co/LFTKRaUJej",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1552815661558116352",
											"sortIndex": "1790198958568505318",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNTUyODE1NjYxNTU4MTE2MzUy",
															"rest_id": "1552815661558116352",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Fri Jul 29 00:38:18 +0000 2022",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Soccer and esports are my life",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1227,
																"followers_count": 19,
																"friends_count": 494,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 0,
																"location": "Venezuela",
																"media_count": 6,
																"name": "Diego Ávila",
																"normal_followers_count": 19,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1552815661558116352/1679557088",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1622871105936424961/q0yaip52_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "davila03lol",
																"statuses_count": 187,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790181701521465344",
											"sortIndex": "1790198958568505317",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTgxNzAxNTIxNDY1MzQ0",
															"rest_id": "1790181701521465344",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue May 14 00:45:57 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "NFT artist ,digital design, 3D artist , NFT creator ,NFT collectionn",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1,
																"followers_count": 0,
																"friends_count": 6,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Kristen Perry",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790181701521465344/1715647606",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790181899043733504/fvXFSXIp_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "KristenPer2379",
																"statuses_count": 1,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1445457530352590862",
											"sortIndex": "1790198958568505316",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNDQ1NDU3NTMwMzUyNTkwODYy",
															"rest_id": "1445457530352590862",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Tue Oct 05 18:35:52 +0000 2021",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 207,
																"followers_count": 112,
																"friends_count": 64,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 1,
																"location": "",
																"media_count": 0,
																"name": "Gerald Allen",
																"normal_followers_count": 112,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1445457730911703056/uayI9OdU_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "GeraldA68046743",
																"statuses_count": 31,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-835839404",
											"sortIndex": "1790198958568505315",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjo4MzU4Mzk0MDQ=",
															"rest_id": "835839404",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Thu Sep 20 15:50:52 +0000 2012",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 35627,
																"followers_count": 164,
																"friends_count": 4874,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2,
																"location": "Guatemala",
																"media_count": 5,
																"name": "Rodrigo",
																"normal_followers_count": 164,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/835839404/1443922508",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/649659213300396032/DiPLKmbs_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "rodrisolorzano7",
																"statuses_count": 298,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-3058315881",
											"sortIndex": "1790198958568505314",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjozMDU4MzE1ODgx",
															"rest_id": "3058315881",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Tue Feb 24 07:46:38 +0000 2015",
																"default_profile": false,
																"default_profile_image": false,
																"description": "I embrace traditional values, modern day Interpretations ,don’t mess with our guns, the smaller the GOVT. the better: truth, justice & the American way #MAGA",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 13488,
																"followers_count": 2013,
																"friends_count": 4606,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 8,
																"location": "New Jersey, USA ",
																"media_count": 295,
																"name": "RoarANastie",
																"normal_followers_count": 2013,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/3058315881/1431505426",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/575545033044062208/0dUyig77_normal.jpeg",
																"profile_interstitial_type": "",
																"screen_name": "RoarANastie",
																"statuses_count": 18741,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-204952097",
											"sortIndex": "1790198958568505313",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoyMDQ5NTIwOTc=",
															"rest_id": "204952097",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": false,
																"created_at": "Tue Oct 19 20:52:41 +0000 2010",
																"default_profile": false,
																"default_profile_image": false,
																"description": "Mamá, esposa y médica ecuatoriana \uD83C\uDDEA\uD83C\uDDE8 Living in Miami for God's purposes \uD83C\uDDFA\uD83C\uDDF8",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 601,
																"followers_count": 435,
																"friends_count": 426,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 3,
																"location": "Coral Gables, FL",
																"media_count": 315,
																"name": "Ely Arosemena \uD83C\uDF3B",
																"normal_followers_count": 435,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/204952097/1706068804",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1748935998094950400/ITLhun66_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "ely_arosemena",
																"statuses_count": 7708,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-2796734758",
											"sortIndex": "1790198958568505312",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoyNzk2NzM0NzU4",
															"rest_id": "2796734758",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Tue Sep 30 23:25:49 +0000 2014",
																"default_profile": false,
																"default_profile_image": false,
																"description": "creator of @engiexchange / ex-eng @apple @kernelco / inventor / gamer / cellist / mountaineer",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "engi.exchange",
																				"expanded_url": "https://engi.exchange",
																				"url": "https://t.co/X82JNi0SY0",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 6242,
																"followers_count": 609,
																"friends_count": 3673,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 31,
																"location": "miami",
																"media_count": 106,
																"name": "garrett maring",
																"normal_followers_count": 609,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/2796734758/1708018793",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1772766722312851456/w3yNop6l_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "enginirmata",
																"statuses_count": 1708,
																"translator_type": "none",
																"url": "https://t.co/X82JNi0SY0",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1715534802940579840",
											"sortIndex": "1790198958568505311",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzE1NTM0ODAyOTQwNTc5ODQw",
															"rest_id": "1715534802940579840",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Sat Oct 21 01:05:59 +0000 2023",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Digital Freight Platform - operating online Logistics worldwide",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 254,
																"followers_count": 100,
																"friends_count": 415,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 30,
																"name": "Ape Global",
																"normal_followers_count": 100,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1715535010499825664/7RDxZ7y0_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "ApeGlobalUSA",
																"statuses_count": 591,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1786207064793944355",
																"professional_type": "Business",
																"category": [
																	{
																		"id": 477,
																		"name": "Professional Services",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790161582015598592",
											"sortIndex": "1790198958568505310",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTYxNTgyMDE1NTk4NTky",
															"rest_id": "1790161582015598592",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 23:26:00 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Degen meme coin trading  \uD83E\uDD21 Low cap hunter",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1,
																"followers_count": 0,
																"friends_count": 6,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": " Alexandria ",
																"media_count": 0,
																"name": "Erica Christy",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790161582015598592/1715642807",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790161757106819072/uO9KJTBx_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "christy_er76072",
																"statuses_count": 1,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1456612360492302341",
											"sortIndex": "1790198958568505309",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNDU2NjEyMzYwNDkyMzAyMzQx",
															"rest_id": "1456612360492302341",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Fri Nov 05 13:20:35 +0000 2021",
																"default_profile": true,
																"default_profile_image": false,
																"description": "21yo artist, self-documenting               PROJECT: BeetleFight!",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "instagram.com/aaron_btlfght/",
																				"expanded_url": "https://www.instagram.com/aaron_btlfght/",
																				"url": "https://t.co/0JKNhHjUOk",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 14936,
																"followers_count": 54,
																"friends_count": 444,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 1,
																"location": "",
																"media_count": 260,
																"name": "aaron",
																"normal_followers_count": 54,
																"pinned_tweet_ids_str": [
																	"1678919039605776385"
																],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1456612360492302341/1713847829",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1782632903316295680/snCzBAF9_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "aaron_btlfght",
																"statuses_count": 544,
																"translator_type": "none",
																"url": "https://t.co/0JKNhHjUOk",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1503566344872247300",
																"professional_type": "Creator",
																"category": [
																	{
																		"id": 1017,
																		"name": "Artist",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790148219214675968",
											"sortIndex": "1790198958568505308",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTQ4MjE5MjE0Njc1OTY4",
															"rest_id": "1790148219214675968",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 22:32:54 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "AI researcher and artist nn",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1,
																"followers_count": 0,
																"friends_count": 5,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Ana Fermin",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790148219214675968/1715639620",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790148393639006208/YT_aXIFg_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "AnaFermin157292",
																"statuses_count": 1,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1731799808619294720",
											"sortIndex": "1790198958568505307",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzMxNzk5ODA4NjE5Mjk0NzIw",
															"rest_id": "1731799808619294720",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon Dec 04 22:17:20 +0000 2023",
																"default_profile": true,
																"default_profile_image": false,
																"description": "International Insurance Reinsurance Exec. & Photographer Capturing the Best in Humans. GFX100S Fuji",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "hugoacardona.com",
																				"expanded_url": "http://www.hugoacardona.com",
																				"url": "https://t.co/nHWfBqrrTx",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1192,
																"followers_count": 113,
																"friends_count": 153,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "Miami, FL",
																"media_count": 18,
																"name": "Hugo A Cardona",
																"normal_followers_count": 113,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1731799808619294720/1702789975",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1736252459129933825/rcQAoaih_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "HACardPhoto",
																"statuses_count": 34,
																"translator_type": "none",
																"url": "https://t.co/nHWfBqrrTx",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-920287604486139904",
											"sortIndex": "1790198958568505306",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjo5MjAyODc2MDQ0ODYxMzk5MDQ=",
															"rest_id": "920287604486139904",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"protected": true,
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Tue Oct 17 13:57:17 +0000 2017",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 913,
																"followers_count": 79,
																"friends_count": 224,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 3,
																"name": "Oris",
																"normal_followers_count": 79,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1788675529157070848/pZSom5_Q_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "oriskemp",
																"statuses_count": 458,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790142637682388992",
											"sortIndex": "1790198958568505305",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTQyNjM3NjgyMzg4OTky",
															"rest_id": "1790142637682388992",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 22:10:49 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Laura Lee",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "LauraLe87177746",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1669320498398781441",
											"sortIndex": "1790198958568505304",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNjY5MzIwNDk4Mzk4NzgxNDQx",
															"rest_id": "1669320498398781441",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": true,
																"created_at": "Thu Jun 15 12:27:11 +0000 2023",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Hello, I’m heidi, I’m very good, fetish, friendly and naughty\uD83D\uDD1E\uD83D\uDE08\uD83D\uDCA6 hmu menu, telegram-#premmbbb, iMessage-#heidimartinz121412@gmail.com",
																"entities": {
																	"description": {
																		"urls": []
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "t.me",
																				"expanded_url": "https://t.me",
																				"url": "https://t.co/s6Huw1vviN",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 65,
																"followers_count": 86,
																"friends_count": 138,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 19,
																"name": "premm_H*️⃣",
																"normal_followers_count": 86,
																"pinned_tweet_ids_str": [
																	"1775593412596465780"
																],
																"possibly_sensitive": true,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1669320498398781441/1712177917",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1775535664143962112/-USU1IQt_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "premmheidi",
																"statuses_count": 67,
																"translator_type": "none",
																"url": "https://t.co/s6Huw1vviN",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"professional": {
																"rest_id": "1775506119688950000",
																"professional_type": "Creator",
																"category": [
																	{
																		"id": 15,
																		"name": "Entertainment & Recreation",
																		"icon_name": "IconBriefcaseStroke"
																	}
																]
															},
															"tipjar_settings": {
																"is_enabled": false,
																"bitcoin_handle": "14CRWzg99NKnKgsQBWokt4kJKaSRXJdk8R",
																"cash_app_handle": ""
															}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790141219739144192",
											"sortIndex": "1790198958568505303",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTQxMjE5NzM5MTQ0MTky",
															"rest_id": "1790141219739144192",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 22:05:10 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Tegan Davison",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "tegan26217",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790140895850823683",
											"sortIndex": "1790198958568505302",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTQwODk1ODUwODIzNjgz",
															"rest_id": "1790140895850823683",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 22:03:52 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Tammy Wells",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "WellsTammy82304",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790139206871732225",
											"sortIndex": "1790198958568505301",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM5MjA2ODcxNzMyMjI1",
															"rest_id": "1790139206871732225",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": true,
																"can_media_tag": false,
																"created_at": "Mon May 13 21:57:25 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Just chillin and joking come vibe \nFollow me on IG@xMunchiiez and\nhttps://t.co/NvDV1heOEp",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "twitch.tv/xMunchiiez",
																				"expanded_url": "http://twitch.tv/xMunchiiez",
																				"url": "https://t.co/NvDV1heOEp",
																				"indices": [
																					66,
																					89
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "youtube.com/@xMunchiiez",
																				"expanded_url": "https://www.youtube.com/@xMunchiiez",
																				"url": "https://t.co/0cMwZhj1ob",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 1,
																"friends_count": 8,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "SzOmbii",
																"normal_followers_count": 1,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790139311104413696/QGFaV33F_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "xMunchiiez",
																"statuses_count": 0,
																"translator_type": "none",
																"url": "https://t.co/0cMwZhj1ob",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790137793831727104",
											"sortIndex": "1790198958568505300",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM3NzkzODMxNzI3MTA0",
															"rest_id": "1790137793831727104",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 21:51:56 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Tracy Spradling",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "SpradlingT91207",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790137385520349184",
											"sortIndex": "1790198958568505299",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM3Mzg1NTIwMzQ5MTg0",
															"rest_id": "1790137385520349184",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 21:50:00 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 6,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Karen Smith",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "KarenSm23441482",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790137215072276480",
											"sortIndex": "1790198958568505298",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM3MjE1MDcyMjc2NDgw",
															"rest_id": "1790137215072276480",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 21:49:21 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Krystal Thurman",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "ThurmanKry26366",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1787013257783341056",
											"sortIndex": "1790198958568505297",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzg3MDEzMjU3NzgzMzQxMDU2",
															"rest_id": "1787013257783341056",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Sun May 05 06:55:47 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 79,
																"followers_count": 2,
																"friends_count": 9,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "mikedawg",
																"normal_followers_count": 2,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "mikedawg0305",
																"statuses_count": 3,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-301672262",
											"sortIndex": "1790198958568505296",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjozMDE2NzIyNjI=",
															"rest_id": "301672262",
															"affiliates_highlighted_label": {},
															"has_graduated_access": true,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Thu May 19 21:18:47 +0000 2011",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 25,
																"followers_count": 11,
																"friends_count": 168,
																"has_custom_timelines": true,
																"is_translator": false,
																"listed_count": 2,
																"location": "",
																"media_count": 0,
																"name": "Alberto",
																"normal_followers_count": 11,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/2170726417/image_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "adelcastillo79",
																"statuses_count": 74,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790134949418250241",
											"sortIndex": "1790198958568505295",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM0OTQ5NDE4MjUwMjQx",
															"rest_id": "1790134949418250241",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"protected": true,
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Mon May 13 21:40:18 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Terri Stiles",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "StilesTerr2186",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790133797641441280",
											"sortIndex": "1790198958568505294",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTMzNzk3NjQxNDQxMjgw",
															"rest_id": "1790133797641441280",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"protected": true,
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Mon May 13 21:35:45 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Isabel Hobbs",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "IsabelHobb13795",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790134469573136384",
											"sortIndex": "1790198958568505293",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTM0NDY5NTczMTM2Mzg0",
															"rest_id": "1790134469573136384",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 21:38:19 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Michelle Smith",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "MSmith7363",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790098747768147968",
											"sortIndex": "1790198958568505292",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMDk4NzQ3NzY4MTQ3OTY4",
															"rest_id": "1790098747768147968",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 19:16:35 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "Life is to short",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 0,
																"friends_count": 38,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "Kano",
																"media_count": 0,
																"name": "Nura Salisu Imam",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_banner_url": "https://pbs.twimg.com/profile_banners/1790098747768147968/1715635471",
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790130424863154176/ojb0GD7F_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "nurahannour31",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790126748270542848",
											"sortIndex": "1790198958568505291",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTI2NzQ4MjcwNTQyODQ4",
															"rest_id": "1790126748270542848",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 21:07:40 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Amy Morgan",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "MorganAmy70117",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790124586677198848",
											"sortIndex": "1790198958568505290",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTI0NTg2Njc3MTk4ODQ4",
															"rest_id": "1790124586677198848",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:59:06 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 4,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Sally Lewis",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "SallyLewis98238",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1784311199674363904",
											"sortIndex": "1790198958568505289",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzg0MzExMTk5Njc0MzYzOTA0",
															"rest_id": "1784311199674363904",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Sat Apr 27 19:58:41 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 131,
																"followers_count": 6,
																"friends_count": 50,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "lil g",
																"normal_followers_count": 6,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1784311488515055616/r9F2w_hn_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "itssscuevas",
																"statuses_count": 14,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790120344457490432",
											"sortIndex": "1790198958568505288",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTIwMzQ0NDU3NDkwNDMy",
															"rest_id": "1790120344457490432",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:44:44 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 0,
																"friends_count": 11,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Jacob Sokolowski",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790121022038949888/BUV_pAwn_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "jacobpski2",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790121635091038208",
											"sortIndex": "1790198958568505287",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTIxNjM1MDkxMDM4MjA4",
															"rest_id": "1790121635091038208",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"protected": true,
																"can_dm": false,
																"can_media_tag": false,
																"created_at": "Mon May 13 20:47:21 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Sarah Martinez",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "SarahMarti59480",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790121135582920705",
											"sortIndex": "1790198958568505286",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTIxMTM1NTgyOTIwNzA1",
															"rest_id": "1790121135582920705",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:45:23 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Kelly Smith",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "KellySmith2073",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790116698877476864",
											"sortIndex": "1790198958568505285",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE2Njk4ODc3NDc2ODY0",
															"rest_id": "1790116698877476864",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:28:08 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Lindsey Smith",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "LindseySmi70790",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790115461662736384",
											"sortIndex": "1790198958568505284",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE1NDYxNjYyNzM2Mzg0",
															"rest_id": "1790115461662736384",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:22:54 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Sheila Wooten",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "SheilaWoot89938",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790116501699047424",
											"sortIndex": "1790198958568505283",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE2NTAxNjk5MDQ3NDI0",
															"rest_id": "1790116501699047424",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:27:01 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Umon Munsen",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "UMunsen91283",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790115887162245120",
											"sortIndex": "1790198958568505282",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE1ODg3MTYyMjQ1MTIw",
															"rest_id": "1790115887162245120",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:24:30 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Kim Hunter",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "KimHunter278166",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790115582173392896",
											"sortIndex": "1790198958568505281",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE1NTgyMTczMzkyODk2",
															"rest_id": "1790115582173392896",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:23:19 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 1,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Megan Jones",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "MeganJones75476",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790116009153605632",
											"sortIndex": "1790198958568505280",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE2MDA5MTUzNjA1NjMy",
															"rest_id": "1790116009153605632",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:25:15 +0000 2024",
																"default_profile": true,
																"default_profile_image": false,
																"description": "For growing your new business or to rank your business or if you need high-quality services, you are in the right place. We provide high.https://t.co/VvJzQJbsqp",
																"entities": {
																	"description": {
																		"urls": [
																			{
																				"display_url": "usasmmit.com",
																				"expanded_url": "https://usasmmit.com",
																				"url": "https://t.co/VvJzQJbsqp",
																				"indices": [
																					137,
																					160
																				]
																			}
																		]
																	},
																	"url": {
																		"urls": [
																			{
																				"display_url": "usasmmit.com/service/buy-ve…",
																				"expanded_url": "https://usasmmit.com/service/buy-verificd-cosh-app-accounts/",
																				"url": "https://t.co/PAXRIQSa8S",
																				"indices": [
																					0,
																					23
																				]
																			}
																		]
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 0,
																"followers_count": 0,
																"friends_count": 5,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "united stets",
																"media_count": 0,
																"name": "Eva Clark",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://pbs.twimg.com/profile_images/1790116459743469568/QfkFXkzF_normal.jpg",
																"profile_interstitial_type": "",
																"screen_name": "EvaClark88137",
																"statuses_count": 0,
																"translator_type": "none",
																"url": "https://t.co/PAXRIQSa8S",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790114561518915586",
											"sortIndex": "1790198958568505279",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTE0NTYxNTE4OTE1NTg2",
															"rest_id": "1790114561518915586",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:19:22 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Danielle Miller",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "DanielleMi94682",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790112799487574017",
											"sortIndex": "1790198958568505278",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTEyNzk5NDg3NTc0MDE3",
															"rest_id": "1790112799487574017",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:12:18 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Michelle Hesano",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "HesanoMich71894",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790110738700267520",
											"sortIndex": "1790198958568505277",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTEwNzM4NzAwMjY3NTIw",
															"rest_id": "1790110738700267520",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 20:04:03 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 4,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Cristina Miller",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "CristinaMi4736",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790107459559837697",
											"sortIndex": "1790198958568505276",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTA3NDU5NTU5ODM3Njk3",
															"rest_id": "1790107459559837697",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 19:51:07 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 2,
																"followers_count": 0,
																"friends_count": 3,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Lauren Smith",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "LaurenSmit71829",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "user-1790108644349980672",
											"sortIndex": "1790198958568505275",
											"content": {
												"entryType": "TimelineTimelineItem",
												"__typename": "TimelineTimelineItem",
												"itemContent": {
													"itemType": "TimelineUser",
													"__typename": "TimelineUser",
													"user_results": {
														"result": {
															"__typename": "User",
															"id": "VXNlcjoxNzkwMTA4NjQ0MzQ5OTgwNjcy",
															"rest_id": "1790108644349980672",
															"affiliates_highlighted_label": {},
															"has_graduated_access": false,
															"is_blue_verified": false,
															"profile_image_shape": "Circle",
															"legacy": {
																"can_dm": false,
																"can_media_tag": true,
																"created_at": "Mon May 13 19:55:45 +0000 2024",
																"default_profile": true,
																"default_profile_image": true,
																"description": "",
																"entities": {
																	"description": {
																		"urls": []
																	}
																},
																"fast_followers_count": 0,
																"favourites_count": 3,
																"followers_count": 0,
																"friends_count": 4,
																"has_custom_timelines": false,
																"is_translator": false,
																"listed_count": 0,
																"location": "",
																"media_count": 0,
																"name": "Ashley Randle",
																"normal_followers_count": 0,
																"pinned_tweet_ids_str": [],
																"possibly_sensitive": false,
																"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
																"profile_interstitial_type": "",
																"screen_name": "ashley_ran98551",
																"statuses_count": 0,
																"translator_type": "none",
																"verified": false,
																"want_retweets": false,
																"withheld_in_countries": []
															},
															"tipjar_settings": {}
														}
													},
													"userDisplayType": "User"
												},
												"clientEventInfo": {
													"component": "FollowersSgs",
													"element": "user"
												}
											}
										},
										{
											"entryId": "cursor-bottom-1790198958568505274",
											"sortIndex": "1790198958568505274",
											"content": {
												"entryType": "TimelineTimelineCursor",
												"__typename": "TimelineTimelineCursor",
												"value": "1798968831783875048|1790198958568505272",
												"cursorType": "Bottom"
											}
										},
										{
											"entryId": "cursor-top-1790198958568505345",
											"sortIndex": "1790198958568505345",
											"content": {
												"entryType": "TimelineTimelineCursor",
												"__typename": "TimelineTimelineCursor",
												"value": "-1|1790198958568505345",
												"cursorType": "Top"
											}
										}
									]
								}
							]
						}
					}
				}
			}
		}
	}`

	legacies, err := parseLegacyInfo(jsonString)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, legacy := range legacies {
		fmt.Printf("Legacy Info: %+v\n", legacy)
	}
}
