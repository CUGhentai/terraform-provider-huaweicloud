---
subcategory: "IoT Device Access (IoTDA)"
---

# huaweicloud_iotda_product

Manages an IoTDA product within HuaweiCloud.

## Example Usage

```hcl
variable "productName" {}

resource "huaweicloud_iotda_space" "space" {
  name = "first_space"
}

resource "huaweicloud_iotda_product" "test" {
  name              = var.productName
  device_type       = "WaterMeter"
  protocol          = "MQTT"
  space_id          = huaweicloud_iotda_space.test.id
  data_type         = "json"
  manufacturer_name = "demo_manufacturer_name"
  industry          = "demo_industry"

  services {
    id   = "service_1"
    type = "serv_type"

    properties {
      name        = "p_1"
      type        = "int"
      min         = "3"
      max         = "666"
      description = "desc"
      method      = "RW"
    }

    properties {
      name       = "p_2"
      type       = "string"
      max_length = 20
      enum_list  = ["1", "E"]
      method     = "R"
    }

    properties {
      name       = "p_3"
      type       = "string"
      method     = "W"
      max_length = 200
    }

    properties {
      name   = "p_4"
      type   = "decimal"
      method = "W"
      min    = "3.1"
      max    = "666.99"
    }

    commands {
      name = "cmd_1"

      paras {
        name        = "cmd_p_1"
        type        = "int"
        description = "desc"
        min         = "1"
        max         = "33"
      }

      responses {
        name        = "cmd_r_1"
        type        = "int"
        description = "desc"
        min         = "1"
        max         = "22"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA product resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the product name. The name contains a maximum of 64 characters. Only letters,
Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are allowed: `?'#().,&%@!`.

* `protocol` - (Required, String) Specifies the protocol.
The valid values are **MQTT**, **CoAP**, **HTTP**, **HTTPS**, **Modbus**, **ONVIF**, **OPC-UA**, **OPC-DA**, **Other**.

* `data_type` - (Required, String) Specifies the type of data.
The valid values are **json** and **binary**.

* `device_type` - (Required, String) Specifies the device type. The device type contains a maximum of 32 characters.
Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters
are allowed: `?'#().,&%@!`. Example: StreetLight, GasMeter, or WaterMeter.

* `services` - (Required, List) Specifies the list of services.
The [services](#IoTDA_service) structure is documented below.

* `manufacturer_name` - (Optional, String) Specifies the manufacturer name. The name contains a maximum of 32 characters.
Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special
characters are allowed: `?'#().,&%@!`.

* `industry` - (Optional, String) Specifies the industry which the device belongs to. The industry contains a maximum of
64 characters. Only letters, Chinese characters, digits, hyphens (-), underscores (_) and
the following special characters are allowed: `?'#().,&%@!`.

* `description` - (Optional, String) Specifies the description of product. The description contains a maximum of 128
characters. Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters
are allowed: `?'#().,&%@!`.

* `product_id` - (Optional, String, ForceNew) Specifies the product ID. The product ID contains a maximum of 32
characters. Only letters, digits, hyphens (-) and underscores (_) are allowed. If omitted, the platform will
automatically allocate a product ID. Changing this parameter will create a new resource.

* `space_id` - (Optional, String, ForceNew) Specifies the resource space ID which the product belongs to. If omitted,
the product will belong to the default resource space. Changing this parameter will create a new resource.

<a name="IoTDA_service"></a>
The `services` block supports:

* `id` - (Required, String) Specifies the service ID. The ID contains a maximum of 64 characters. Only letters,
Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are allowed: `?'#().,&%@!`.

* `type` - (Optional, String) Specifies the service type. The type contains a maximum of 64 characters. Only letters,
Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are allowed: `?'#().,&%@!`.
The default value is equal to service ID.

* `description` - (Optional, String) Specifies description of service. The description contains a maximum of 128
characters. Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special
characters are allowed: `?'#().,&%@!`.

* `properties` - (Optional, List) Specifies the list of properties for the service.
The [properties](#IoTDA_service_properties) structure is documented below.

* `commands` - (Optional, List) Specifies the list of commands for the service.
The [commands](#IoTDA_service_commands) structure is documented below.

<a name="IoTDA_service_properties"></a>
The `properties` block supports:

* `name` - (Required, String) Specifies the name of the parameter. The name contains a maximum of 64 characters.
Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are
allowed: `?'#().,&%@!`.

* `type` - (Required, String) Specifies the type of the parameter.
The valid values are **int**, **decimal**, **string**, **DateTime**, **jsonObject** and **string list**.

* `access_permission` - (Required, String) Specifies the access permission of the parameter. The options are as follows:
  + **R**: readable.
  + **w**: writable.
  + **RW**: readable and writable.

* `description` - (Optional, String) Specifies the description of the parameter. The description contains a maximum of
128 characters. Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special
characters are allowed: `?'#().,&%@!`.

* `min` - (Optional, String) Specifies the min value of the parameter when the `type` is **int** or **decimal**.
Value range: -2147483647 ~ 2147483647. Defaults to **"0"**.

* `max` - (Optional, String) Specifies the max value of the parameter when the `type` is **int** or **decimal**.
Value range: -2147483647 ~ 2147483647. Defaults to **"65535"**.

* `step` - (Optional, Float) Specifies the step of the parameter when the `type` is **int** or **decimal**.
Value range: 0 ~ 2147483647. Defaults to **0**.

* `unit` - (Optional, String) Specifies the unit of the parameter when the `type` is **int** or **decimal**.
The unit contains a maximum of 16 characters.

* `max_length` - (Optional, Int) Specifies the max length of the parameter when the `type` is **string**, **DateTime**,
**jsonObject** or **string list**. Value range: 0 ~ 2147483647. Defaults to **0**.

<a name="IoTDA_service_commands"></a>
The `commands` block supports:

* `name` - (Required, String) Specifies the name of the command. The name contains a maximum of 64 characters.
Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters
are allowed: `?'#().,&%@!`.

* `paras` - (Optional, List) Specifies the list of parameters for the command.
The [paras](#IoTDA_service_commands_properties) structure is documented below.

* `responses` - (Optional, List) Specifies the list of responses for the command.
The [responses](#IoTDA_service_commands_properties) structure is documented below.

<a name="IoTDA_service_commands_properties"></a>
The `paras` and `responses` block supports:

* `name` - (Required, String) Specifies the name of the parameter. The name contains a maximum of 64 characters.
Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special characters are
allowed: `?'#().,&%@!`.

* `type` - (Optional, String) Specifies the type of the parameter.
The valid values are **int**, **decimal**, **string**, **DateTime**, **jsonObject** and **string list**.

* `description` - (Optional, String) Specifies the description of the parameter. The description contains a maximum of
128 characters. Only letters, Chinese characters, digits, hyphens (-), underscores (_) and the following special
characters are allowed: `?'#().,&%@!`.

* `min` - (Optional, String) Specifies the min value of the parameter when the `type` is **int** or **decimal**.
Value range: -2147483647 ~ 2147483647. Defaults to **"0"**.

* `max` - (Optional, String) Specifies the max value of the parameter when the `type` is **int** or **decimal**.
Value range: -2147483647 ~ 2147483647. Defaults to **"65535"**.

* `step` - (Optional, Float) Specifies the step of the parameter when the `type` is **int** or **decimal**.
Value range: 0 ~ 2147483647. Defaults to **0**.

* `unit` - (Optional, String) Specifies the unit of the parameter when the `type` is **int** or **decimal**.
The unit contains a maximum of 16 characters.

* `max_length` - (Optional, Int) Specifies the max length of the parameter when the `type` is **string**, **DateTime**,
**jsonObject** or **string list**. Value range: 0 ~ 2147483647. Defaults to **0**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Products can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iotda_product.test 10022532f4f94f26b01daa1e424853e1
```
