package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mackerelio/mackerel-client-go"
)

func resourceMackerelServiceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMackerelServiceMonitorCreate,
		Read:   resourceMackerelServiceMonitorRead,
		Update: resourceMackerelServiceMonitorUpdate,
		Delete: resourceMackerelServiceMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"warning": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"critical": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"notification_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_mute": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceMackerelServiceMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                 "service",
		Name:                 d.Get("name").(string),
		Service:              d.Get("service").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	monitor, err := client.CreateMonitor(input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q created.", monitor.MonitorID())
	d.SetId(monitor.MonitorID())

	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	log.Printf("[DEBUG] Reading mackerel monitor: %q", d.Id())
	monitors, err := client.FindMonitors()
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		if monitor.MonitorType() == "service" && monitor.MonitorID() == d.Id() {
			mon := monitor.(*mackerel.MonitorServiceMetric)
			_ = d.Set("id", mon.ID)
			_ = d.Set("name", mon.Name)
			_ = d.Set("service", mon.Service)
			_ = d.Set("duration", mon.Duration)
			_ = d.Set("metric", mon.Metric)
			_ = d.Set("operator", mon.Operator)
			_ = d.Set("warning", mon.Warning)
			_ = d.Set("critical", mon.Critical)
			_ = d.Set("notification_interval", mon.NotificationInterval)
			_ = d.Set("is_mute", mon.IsMute)
			break
		}
	}

	return nil
}

func resourceMackerelServiceMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	input := &mackerel.MonitorServiceMetric{
		Type:                 "service",
		Name:                 d.Get("name").(string),
		Service:              d.Get("service").(string),
		Duration:             uint64(d.Get("duration").(int)),
		Metric:               d.Get("metric").(string),
		Operator:             d.Get("operator").(string),
		Warning:              pfloat64(d.Get("warning").(float64)),
		Critical:             pfloat64(d.Get("critical").(float64)),
		NotificationInterval: uint64(d.Get("notification_interval").(int)),
		IsMute:               d.Get("is_mute").(bool),
	}

	_, err := client.UpdateMonitor(d.Id(), input)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q updated.", d.Id())
	return resourceMackerelServiceMonitorRead(d, meta)
}

func resourceMackerelServiceMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mackerel.Client)

	_, err := client.DeleteMonitor(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] mackerel monitor %q deleted.", d.Id())
	d.SetId("")

	return nil
}
