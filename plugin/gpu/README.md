Example output from `nvidia-smi -q -x`

```xml

<?xml version="1.0" ?>
<!DOCTYPE nvidia_smi_log SYSTEM "nvsmi_device_v10.dtd">
<nvidia_smi_log>
	<timestamp>Tue Oct 29 16:27:25 2019</timestamp>
	<driver_version>418.56</driver_version>
	<cuda_version>10.1</cuda_version>
	<attached_gpus>1</attached_gpus>
	<gpu id="00000000:01:00.0">
		<product_name>GeForce GTX 1070</product_name>
		<product_brand>GeForce</product_brand>
		<display_mode>Enabled</display_mode>
		<display_active>Enabled</display_active>
		<persistence_mode>Enabled</persistence_mode>
		<accounting_mode>Disabled</accounting_mode>
		<accounting_mode_buffer_size>4000</accounting_mode_buffer_size>
		<driver_model>
			<current_dm>N/A</current_dm>
			<pending_dm>N/A</pending_dm>
		</driver_model>
		<serial>N/A</serial>
		<uuid>GPU-aaaaaaaa-bbbb-cccc-dddd-eeeeeeffffff</uuid>
		<minor_number>0</minor_number>
		<vbios_version>11.22.33.44.55</vbios_version>
		<multigpu_board>No</multigpu_board>
		<board_id>0x100</board_id>
		<gpu_part_number>N/A</gpu_part_number>
		<inforom_version>
			<img_version>G001.0000.00.00</img_version>
			<oem_object>1.1</oem_object>
			<ecc_object>N/A</ecc_object>
			<pwr_object>N/A</pwr_object>
		</inforom_version>
		<gpu_operation_mode>
			<current_gom>N/A</current_gom>
			<pending_gom>N/A</pending_gom>
		</gpu_operation_mode>
		<gpu_virtualization_mode>
			<virtualization_mode>None</virtualization_mode>
		</gpu_virtualization_mode>
		<ibmnpu>
			<relaxed_ordering_mode>N/A</relaxed_ordering_mode>
		</ibmnpu>
		<pci>
			<pci_bus>01</pci_bus>
			<pci_device>00</pci_device>
			<pci_domain>0000</pci_domain>
			<pci_device_id>1B11111B</pci_device_id>
			<pci_bus_id>00000000:01:00.0</pci_bus_id>
			<pci_sub_system_id>11111111</pci_sub_system_id>
			<pci_gpu_link_info>
				<pcie_gen>
					<max_link_gen>3</max_link_gen>
					<current_link_gen>1</current_link_gen>
				</pcie_gen>
				<link_widths>
					<max_link_width>16x</max_link_width>
					<current_link_width>16x</current_link_width>
				</link_widths>
			</pci_gpu_link_info>
			<pci_bridge_chip>
				<bridge_chip_type>N/A</bridge_chip_type>
				<bridge_chip_fw>N/A</bridge_chip_fw>
			</pci_bridge_chip>
			<replay_counter>0</replay_counter>
			<replay_rollover_counter>0</replay_rollover_counter>
			<tx_util>21000 KB/s</tx_util>
			<rx_util>24000 KB/s</rx_util>
		</pci>
		<fan_speed>0 %</fan_speed>
		<performance_state>P8</performance_state>
		<clocks_throttle_reasons>
			<clocks_throttle_reason_gpu_idle>Active</clocks_throttle_reason_gpu_idle>
			<clocks_throttle_reason_applications_clocks_setting>Not Active</clocks_throttle_reason_applications_clocks_setting>
			<clocks_throttle_reason_sw_power_cap>Not Active</clocks_throttle_reason_sw_power_cap>
			<clocks_throttle_reason_hw_slowdown>Not Active</clocks_throttle_reason_hw_slowdown>
			<clocks_throttle_reason_hw_thermal_slowdown>Not Active</clocks_throttle_reason_hw_thermal_slowdown>
			<clocks_throttle_reason_hw_power_brake_slowdown>Not Active</clocks_throttle_reason_hw_power_brake_slowdown>
			<clocks_throttle_reason_sync_boost>Not Active</clocks_throttle_reason_sync_boost>
			<clocks_throttle_reason_sw_thermal_slowdown>Not Active</clocks_throttle_reason_sw_thermal_slowdown>
			<clocks_throttle_reason_display_clocks_setting>Not Active</clocks_throttle_reason_display_clocks_setting>
		</clocks_throttle_reasons>
		<fb_memory_usage>
			<total>8119 MiB</total>
			<used>561 MiB</used>
			<free>7558 MiB</free>
		</fb_memory_usage>
		<bar1_memory_usage>
			<total>256 MiB</total>
			<used>6 MiB</used>
			<free>250 MiB</free>
		</bar1_memory_usage>
		<compute_mode>Default</compute_mode>
		<utilization>
			<gpu_util>40 %</gpu_util>
			<memory_util>29 %</memory_util>
			<encoder_util>0 %</encoder_util>
			<decoder_util>0 %</decoder_util>
		</utilization>
		<encoder_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</encoder_stats>
		<fbc_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</fbc_stats>
		<ecc_mode>
			<current_ecc>N/A</current_ecc>
			<pending_ecc>N/A</pending_ecc>
		</ecc_mode>
		<ecc_errors>
			<volatile>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</volatile>
			<aggregate>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</aggregate>
		</ecc_errors>
		<retired_pages>
			<multiple_single_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</multiple_single_bit_retirement>
			<double_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</double_bit_retirement>
			<pending_retirement>N/A</pending_retirement>
		</retired_pages>
		<temperature>
			<gpu_temp>49 C</gpu_temp>
			<gpu_temp_max_threshold>99 C</gpu_temp_max_threshold>
			<gpu_temp_slow_threshold>96 C</gpu_temp_slow_threshold>
			<gpu_temp_max_gpu_threshold>N/A</gpu_temp_max_gpu_threshold>
			<memory_temp>N/A</memory_temp>
			<gpu_temp_max_mem_threshold>N/A</gpu_temp_max_mem_threshold>
		</temperature>
		<power_readings>
			<power_state>P8</power_state>
			<power_management>Supported</power_management>
			<power_draw>11.41 W</power_draw>
			<power_limit>151.00 W</power_limit>
			<default_power_limit>151.00 W</default_power_limit>
			<enforced_power_limit>151.00 W</enforced_power_limit>
			<min_power_limit>75.00 W</min_power_limit>
			<max_power_limit>170.00 W</max_power_limit>
		</power_readings>
		<clocks>
			<graphics_clock>480 MHz</graphics_clock>
			<sm_clock>480 MHz</sm_clock>
			<mem_clock>405 MHz</mem_clock>
			<video_clock>544 MHz</video_clock>
		</clocks>
		<applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</applications_clocks>
		<default_applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</default_applications_clocks>
		<max_clocks>
			<graphics_clock>1999 MHz</graphics_clock>
			<sm_clock>1999 MHz</sm_clock>
			<mem_clock>4004 MHz</mem_clock>
			<video_clock>1708 MHz</video_clock>
		</max_clocks>
		<max_customer_boost_clocks>
			<graphics_clock>N/A</graphics_clock>
		</max_customer_boost_clocks>
		<clock_policy>
			<auto_boost>N/A</auto_boost>
			<auto_boost_default>N/A</auto_boost_default>
		</clock_policy>
		<supported_clocks>N/A</supported_clocks>
		<processes>
			<process_info>
				<pid>1860</pid>
				<type>G</type>
				<process_name>/usr/lib/xorg/Xorg</process_name>
				<used_memory>343 MiB</used_memory>
			</process_info>
			<process_info>
				<pid>3498</pid>
				<type>G</type>
				<process_name>/usr/share/spotify/spotify --log-file=/usr/share/spotify/debug.log</process_name>
				<used_memory>59 MiB</used_memory>
			</process_info>
		</processes>
		<accounted_processes>
		</accounted_processes>
	</gpu>

</nvidia_smi_log>

```