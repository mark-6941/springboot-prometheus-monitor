
# 資料中心
data "vsphere_datacenter" "dc" {
  name = var.dc
}

# 儲存與資源池
data "vsphere_datastore_cluster" "datastore_cluster" {
  count = var.datastore_cluster != "" ? 1 : 0
  name  = var.datastore_cluster
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "template" {
  count = var.content_library == null ? 1 : 0
  name  = var.vmtemp
  datacenter_id = data.vsphere_datacenter.dc.id
}

# 虛擬機資源
resource "vsphere_virtual_machine" "vm" {
  count = var.instances
  name  = var.staticvmname != null ? var.staticvmname : format("%s-%d", var.vmname, count.index)
  datacenter_id = data.vsphere_datacenter.dc.id

  # 硬體配置
  cpu_count = var.cpu_count
  memory_num_bytes = var.memory_size * 1024 * 1024

  # 網路
  network_interface {
    device_name = var.network
    guest_ip_address = var.ipv4_address
  }

  # 磁碟
  disk {
    label = "virtual-vm-disk"
    size  = var.disk_size
    thin_provisioned = var.thin_provisioned
  }

  depends_on = [var.vm_depends_on]
}

resource "vsphere_virtual_machine" "vm" {
  name             = "app-db-server"
  resource_pool_id = data.vsphere_resource_pool.pool.id
  datastore_id     = data.vsphere_datastore.datastore.id

  num_cpus = 2
  memory   = 4096
  guest_id = "ubuntu64Guest"

  network_interface {
    network_id = data.vsphere_network.network.id
  }

  disk {
    label = "disk0"
    size  = 50
  }
  
  # 部署完成後觸發 Ansible
  provisioner "local-exec" {
    command = "ansible-playbook -i '${self.default_ip_address},' setup-monitoring.yml"
  }
}
