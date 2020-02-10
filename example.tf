provider "multiverse" {}

resource "multiverse_custom_resource" "test1" {
    executor = "node"
    script = "example.js"
    id_key = "ID"
    data = <<JSON
{
    "foo": "bar2",
    "bar": "baz"
}
JSON
}

output "test" {
    value = multiverse_custom_resource.test1
}