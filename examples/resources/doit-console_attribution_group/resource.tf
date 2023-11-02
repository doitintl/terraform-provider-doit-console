# Manage Attribution group
resource "doit-console_attribution_group" "attributeGroup" {
  name         = "attritestnewgroup"
  description  = "attritestgroup"
  attributions = [doit_attribution.attribute1.id, doit_attribution.attribute2.id]
}