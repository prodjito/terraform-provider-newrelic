package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

var (
	validIconValues []string = []string{
		"none",
		"archive",
		"bar-chart",
		"line-chart",
		"bullseye",
		"user",
		"usd",
		"money",
		"thumbs-up",
		"thumbs-down",
		"cloud",
		"bell",
		"bullhorn",
		"comments-o",
		"envelope",
		"globe",
		"shopping-cart",
		"sitemap",
		"clock-o",
		"crosshairs",
		"rocket",
		"users",
		"mobile",
		"tablet",
		"adjust",
		"dashboard",
		"flag",
		"flask",
		"road",
		"bolt",
		"cog",
		"leaf",
		"magic",
		"puzzle-piece",
		"bug",
		"fire",
		"legal",
		"trophy",
		"pie-chart",
		"sliders",
		"paper-plane",
		"life-ring",
		"heart",
	}

	validWidgetVisualizationValues []string = []string{
		"billboard",
		"gauge",
		"billboard_comparison",
		"facet_bar_chart",
		"faceted_line_chart",
		"facet_pie_chart",
		"facet_table",
		"faceted_area_chart",
		"heatmap",
		"attribute_sheet",
		"single_event",
		"histogram",
		"funnel",
		"raw_json",
		"event_feed",
		"event_table",
		"uniques_list",
		"line_chart",
		"comparison_line_chart",
		"markdown",
		"metric_line_chart",
	}
)

func resourceNewRelicDb() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicDbCreate,
		Read:   resourceNewRelicDbRead,
		Update: resourceNewRelicDbUpdate,
		Delete: resourceNewRelicDbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "description": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"icon": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bar-chart",
				ValidateFunc: validation.StringInSlice(validIconValues, false),
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validation.StringInSlice([]string{"owner", "all"}, false),
			},
			"editable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "editable_by_all",
				ValidateFunc: validation.StringInSlice([]string{"read_only", "editable_by_owner", "editable_by_all", "all"}, false),
			},
			"ui_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_types": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"attributes": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"widget": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 300,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"visualization": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(validWidgetVisualizationValues, false),
						},
						"widget_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"account_id": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"data": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nrql": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"duration": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"end_time": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"raw_metric_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"facet": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"order_by": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      0,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"entity_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"metrics": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"units": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"scope": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"values": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"compare_with": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"offset_duration": {
													Type:     schema.TypeString,
													Required: true,
												},
												"presentation": {
													Type:     schema.TypeMap,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"presentation": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": {
										Type:     schema.TypeString,
										Required: true,
									},
									"notes": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"red": {
													Type:     schema.TypeFloat,
													Required: true,
												},
												"yellow": {
													Type:     schema.TypeFloat,
													Optional: true,
													Default:  0,
												},
											},
										},
									},
									"drilldown_dashboard_id": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"layout": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"width": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"height": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"row": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"column": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicDbRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	dashboardID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("title", dashboard.Title)
	d.Set("icon", dashboard.Icon)
	//d.Set("description", dashboard.Description)
	d.Set("visibility", dashboard.Visibility)
	d.Set("editable", dashboard.Editable)
	d.Set("ui_url", dashboard.UIURL)
	d.Set("api_url", dashboard.APIRL)
	d.Set("owner_email", dashboard.OwnerEmail)

	if err = d.Set("filter", flattenDbFilter(&dashboard.Filter)); err != nil {
		return err
	}

	if dashboard.Widgets != nil && len(dashboard.Widgets) > 0 {
		if err = d.Set("widget", flattenWidgets(dashboard.Widgets)); err != nil {
			if err == newrelic.ErrNotFound {
				d.SetId("")
				return nil
			}

			return err
		}
	}

	d.SetId(strconv.Itoa(dashboard.ID))

	return nil
}

func resourceNewRelicDbCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	dashboard, err := expandDb(d)
	if err != nil {
		return err
	}

	dashboard, err = client.CreateDashboard(*dashboard)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(dashboard.ID))

	return resourceNewRelicDbRead(d, meta)
}

func resourceNewRelicDbUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("updating..")
	client := meta.(*ProviderConfig).Client

	dashboard, err := expandDb(d)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard.ID = id

	_, err = client.UpdateDashboard(*dashboard)
	if err != nil {
		return err
	}

	return resourceNewRelicDbRead(d, meta)
}

func resourceNewRelicDbDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteDashboard(id); err != nil {
		return err
	}

	return nil
}
