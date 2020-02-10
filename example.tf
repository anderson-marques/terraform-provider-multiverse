provider "multiverse" {}

resource "multiverse_custom_resource" "test_node" {
    executor = "node"
    script = "example.js"
    id_key = "ID"
    data = <<JSON
{
    "foo": "bar0",
    "bar": "baz"
}
JSON
}

resource "multiverse_custom_resource" "test_bash" {
    executor = "bash"
    script = "example.sh"
    id_key = "ID"
    data = <<JSON
{
    "foo": "bar11",
    "bar": "baz"
}
JSON
}

output "test_node" {
    value = multiverse_custom_resource.test_node
}

output "test_node_deep_json" {
    value = jsondecode(multiverse_custom_resource.test_node.deep_object)
}

output "test_bash" {
    value = multiverse_custom_resource.test_bash
}

output "test_bash_deep_json" {
    value = jsondecode(multiverse_custom_resource.test_bash.deep_object)
}