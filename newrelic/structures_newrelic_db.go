package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func flattenDbFilter(in *newrelic.DashboardFilter) []interface{} {
	if in == nil {
		return nil
	}

	if len(in.Attributes) == 0 && len(in.EventTypes) == 0 {
		return nil
	}

	m := make(map[string]interface{})

	if in.EventTypes != nil && len(in.EventTypes) > 0 {
		m["event_types"] = in.EventTypes
	}

	if in.Attributes != nil && len(in.EventTypes) > 0 {
		m["attributes"] = in.EventTypes
	}

	return []interface{}{m}
}

func flattenWidgets(in []newrelic.DashboardWidget) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, w := range in {
		m := make(map[string]interface{})
		m["widget_id"] = w.ID
		m["visualization"] = w.Visualization
		m["account_id"] = w.AccountID
		m["presentation"] = flattenWidgetPresentation(&w.Presentation)
		m["layout"] = flattenWidgetLayout(&w.Layout)

		if w.Data != nil {
			m["data"] = flattenWidgetData(&w.Data[0])
		}

		out[i] = m
	}

	return out
}

func flattenWidgetPresentation(in *newrelic.DashboardWidgetPresentation) []interface{} {
	m := make(map[string]interface{})

	m["title"] = in.Title

	if in.Notes != "" {
		m["notes"] = in.Notes
	}

	if in.DrilldownDashboardID > 0 {
		m["drilldown_dashboard_id"] = in.DrilldownDashboardID
	}

	if in.Threshold != nil {
		m["threshold"] = flattenWidgetPresentationThreshold(in.Threshold)
	}

	return []interface{}{m}
}

func flattenWidgetPresentationThreshold(in *newrelic.DashboardWidgetThreshold) []interface{} {
	m := make(map[string]interface{})

	m["red"] = in.Red

	if in.Yellow > 0 {
		m["yellow"] = in.Yellow
	}

	return []interface{}{m}
}

func flattenWidgetLayout(in *newrelic.DashboardWidgetLayout) []interface{} {
	m := make(map[string]interface{})

	m["row"] = in.Row
	m["column"] = in.Column
	m["width"] = in.Width
	m["height"] = in.Height

	return []interface{}{m}
}

func flattenWidgetData(in *newrelic.DashboardWidgetData) []interface{} {
	m := make(map[string]interface{})

	if in.NRQL != "" {
		m["nrql"] = in.NRQL
	}

	if in.Source != "" {
		m["source"] = in.Source
	}

	if in.Duration > 0 {
		m["duration"] = in.Duration
	}

	if in.EndTime > 0 {
		m["end_time"] = in.EndTime
	}

	if in.RawMetricName != "" {
		m["raw_metric_name"] = in.RawMetricName
	}

	if in.Facet != "" {
		m["facet"] = in.Facet
	}

	if in.OrderBy != "" {
		m["order_by"] = in.OrderBy
	}

	if in.Limit > 0 {
		m["limit"] = in.Limit
	}

	if in.EntityIds != nil && len(in.EntityIds) > 0 {
		m["entity_ids"] = in.EntityIds
	}

	if in.CompareWith != nil && len(in.CompareWith) > 0 {
		m["compare_with"] = flattenWidgetDataCompareWith(in.CompareWith)
	}

	if in.Metrics != nil && len(in.Metrics) > 0 {
		m["metrics"] = flattenWidgetDataMetrics(in.Metrics)
	}

	return []interface{}{m}
}

func flattenWidgetDataCompareWith(in []newrelic.DashboardWidgetDataCompareWith) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		m["offset_duration"] = v.OffsetDuration
		m["presentation"] = flattenWidgetDataCompareWithPresentation(&v.Presentation)

		out[i] = m
	}

	return out
}

func flattenWidgetDataCompareWithPresentation(in *newrelic.DashboardWidgetDataCompareWithPresentation) interface{} {
	m := make(map[string]interface{})

	m["name"] = in.Name
	m["color"] = in.Color

	return m
}

func flattenWidgetDataMetrics(in []newrelic.DashboardWidgetDataMetric) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		m["name"] = v.Name
		m["units"] = v.Units
		m["scope"] = v.Scope

		if v.Values != nil && len(v.Values) > 0 {
			m["values"] = v.Values
		}

		out[i] = m
	}

	return out
}

func expandDb(d *schema.ResourceData) (*newrelic.Dashboard, error) {
	log.Printf("expanding db..")
	dashboard := &newrelic.Dashboard{
		Title:      d.Get("title").(string),
		Metadata:   newrelic.DashboardMetadata{Version: 1},
		Icon:       d.Get("icon").(string),
		Visibility: d.Get("visibility").(string),
		Editable:   d.Get("editable").(string),
	}

	if filter, ok := d.GetOk("filter"); ok && len(filter.([]interface{})) > 0 {
		dashboard.Filter = expandFilter(filter.([]interface{})[0].(map[string]interface{}))
	}

	if widgets, ok := d.GetOk("widget"); ok && widgets.(*schema.Set).Len() > 0 {
		expandedWidgets, err := expandWidgets(widgets.(*schema.Set).List())

		if err != nil {
			return nil, err
		}

		dashboard.Widgets = expandedWidgets
	}

	return dashboard, nil
}

func expandFilter(filter map[string]interface{}) newrelic.DashboardFilter {
	out := newrelic.DashboardFilter{}

	if v, ok := filter["attributes"]; ok {
		attributes := v.(*schema.Set).List()
		vs := make([]string, 0, len(attributes))
		for _, a := range attributes {
			vs = append(vs, a.(string))
		}

		out.Attributes = vs
	}

	if v, ok := filter["event_types"]; ok {
		eventTypes := v.(*schema.Set).List()
		vs := make([]string, 0, len(eventTypes))
		for _, e := range eventTypes {
			vs = append(vs, e.(string))
		}
		out.EventTypes = vs
	}

	return out
}

func expandWidgets(widgets []interface{}) ([]newrelic.DashboardWidget, error) {
	log.Printf("expanding widgets..")
	if len(widgets) < 1 {
		return []newrelic.DashboardWidget{}, nil
	}

	perms := make([]newrelic.DashboardWidget, len(widgets))

	for i, rawCfg := range widgets {
		cfg := rawCfg.(map[string]interface{})
		expandedWidget, err := expandWidget(cfg)

		if err != nil {
			return nil, err
		}

		perms[i] = *expandedWidget
	}

	return perms, nil
}

func expandWidget(cfg map[string]interface{}) (*newrelic.DashboardWidget, error) {
	log.Printf("expanding widget...")

	widget := &newrelic.DashboardWidget{
		Visualization: cfg["visualization"].(string),
		AccountID:     cfg["account_id"].(int),
	}

	if v, ok := cfg["data"]; ok {
		data := v.([]interface{})[0].(map[string]interface{})
		err := validateWidgetData(cfg)
		if err != nil {
			return nil, err
		}

		widget.Data = expandWidgetData(data)
	}

	if presentation, ok := cfg["presentation"]; ok {
		widget.Presentation = expandWidgetPresentation(presentation.([]interface{})[0].(map[string]interface{}))
	}

	if layout, ok := cfg["layout"]; ok {
		expandedLayout, err := expandWidgetLayout(layout.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		widget.Layout = *expandedLayout
	}

	return widget, nil
}

func validateWidgetData(cfg map[string]interface{}) error {
	visualization := cfg["visualization"].(string)
	log.Printf("Validating!!!")

	data := cfg["data"].([]interface{})[0].(map[string]interface{})
	presentation := cfg["presentation"].([]interface{})[0].(map[string]interface{})

	switch visualization {
	case "billboard", "gauge", "billboard_comparison":
		if _, ok := data["nrql"]; !ok {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}

		if t, ok := presentation["threshold"]; !ok || len(t.([]interface{})) == 0 {
			return fmt.Errorf("threshold is required for %s visualization", visualization)
		}
	case "facet_bar_chart", "faceted_line_chart", "facet_pie_chart", "facet_table", "faceted_area_chart", "heatmap":
		if _, ok := data["nrql"]; !ok {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
	case "attribute_sheet", "single_event", "histogram", "funnel", "raw_json", "event_feed", "event_table", "uniques_list", "line_chart", "comparison_line_chart":
		if _, ok := data["nrql"]; !ok {
			return fmt.Errorf("nrql is required for %s visualization", visualization)
		}
	case "markdown":
		if _, ok := data["source"]; !ok {
			return fmt.Errorf("source is required for %s visualization", visualization)
		}
	case "metric_line_chart":
		if _, ok := data["metrics"]; !ok {
			return fmt.Errorf("metrics is required for %s visualization", visualization)
		}

		if _, ok := data["entity_ids"]; !ok {
			return fmt.Errorf("entity_ids is required for %s visualization", visualization)
		}
	}

	return nil
}

func expandWidgetData(cfg map[string]interface{}) []newrelic.DashboardWidgetData {
	widgetData := newrelic.DashboardWidgetData{}

	if nrql, ok := cfg["nrql"]; ok {
		widgetData.NRQL = nrql.(string)
	}

	if source, ok := cfg["source"]; ok {
		widgetData.Source = source.(string)
	}

	if duration, ok := cfg["duration"]; ok {
		widgetData.Duration = duration.(int)
	}

	if endTime, ok := cfg["end_time"]; ok {
		widgetData.EndTime = endTime.(int)
	}

	if rawMetricName, ok := cfg["raw_metric_name"]; ok {
		widgetData.RawMetricName = rawMetricName.(string)
	}

	if facet, ok := cfg["facet"]; ok {
		widgetData.Facet = facet.(string)
	}

	if orderBy, ok := cfg["order_by"]; ok {
		widgetData.OrderBy = orderBy.(string)
	}

	if limit, ok := cfg["limit"]; ok {
		widgetData.Limit = limit.(int)
	}

	if entityIds, ok := cfg["entity_ids"]; ok && entityIds.(*schema.Set).Len() > 0 {
		widgetData.EntityIds = expandIntSet(entityIds.(*schema.Set))
	}

	if compareWith, ok := cfg["compare_with"]; ok {
		widgetData.CompareWith = expandWidgetDataCompareWith(compareWith.(*schema.Set))
	}

	// widget data is a slice for legacy reasons
	return []newrelic.DashboardWidgetData{widgetData}
}

func expandWidgetDataCompareWith(windows *schema.Set) []newrelic.DashboardWidgetDataCompareWith {
	if windows.Len() < 1 {
		return []newrelic.DashboardWidgetDataCompareWith{}
	}

	perms := make([]newrelic.DashboardWidgetDataCompareWith, windows.Len())

	for i, rawCfg := range windows.List() {
		cfg := rawCfg.(map[string]interface{})

		perms[i] = newrelic.DashboardWidgetDataCompareWith{
			OffsetDuration: cfg["offset_duration"].(string),
			Presentation:   expandWidgetDataCompareWithPresentation(cfg["presentation"].(map[string]interface{})),
		}
	}

	return perms
}

func expandWidgetDataCompareWithPresentation(cfg map[string]interface{}) newrelic.DashboardWidgetDataCompareWithPresentation {

	widgetDataCompareWithPresentation := newrelic.DashboardWidgetDataCompareWithPresentation{
		Name:  cfg["name"].(string),
		Color: cfg["color"].(string),
	}

	return widgetDataCompareWithPresentation
}

func expandWidgetPresentation(cfg map[string]interface{}) newrelic.DashboardWidgetPresentation {
	widgetPresentation := newrelic.DashboardWidgetPresentation{
		Title: cfg["title"].(string),
		Notes: cfg["notes"].(string),
	}

	if threshold, ok := cfg["threshold"]; ok && len(threshold.([]interface{})) > 0 {
		widgetPresentation.Threshold = expandWidgetThreshold(threshold.([]interface{})[0].(map[string]interface{}))
	}

	return widgetPresentation
}

func expandWidgetThreshold(cfg map[string]interface{}) *newrelic.DashboardWidgetThreshold {
	widgetThreshold := &newrelic.DashboardWidgetThreshold{}

	if red, ok := cfg["red"]; ok {
		widgetThreshold.Red = red.(float64)
	}

	if yellow, ok := cfg["yellow"]; ok {
		widgetThreshold.Yellow = yellow.(float64)
	}

	return widgetThreshold
}

func expandWidgetLayout(cfg map[string]interface{}) (*newrelic.DashboardWidgetLayout, error) {
	widgetLayout := &newrelic.DashboardWidgetLayout{
		Row:    cfg["row"].(int),
		Column: cfg["column"].(int),
		Width:  cfg["width"].(int),
		Height: cfg["height"].(int),
	}

	return widgetLayout, nil
}
