<template>
<div class="modal inmodal fade" id="addTask" tabindex="-1" role="dialog"  aria-hidden="true">
    <div class="modal-dialog modal-md">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span class="si-close modal-close" aria-hidden="true"></span><span class="sr-only">Close</span></button>
                <h4 class="modal-title">Create task</h4>
            </div>
            <div class="modal-body">
                <h4>Basic</h4>
                <div class="form-group">
                    <label for="name-field" class="control-label">Task Name</label>
                    <div><input type="text" id="name-field" name="name" class="form-control" v-model="task.name" required="required"></div>
                </div>
                <div class="row">
                    <div class="col-sm-4">
                        <div class="form-group">
                            <label for="cpus-field" class="control-label">Cpus</label>
                            <div><input type="text" id="cpus-field" name="cpus" class="form-control" value="{{form.cpus}}" v-model="task.cpus"></div>
                        </div>                        
                    </div>
                    <div class="col-sm-4">
                        <div class="form-group">
                            <label for="mem-field" class="control-label">Memory(MB)</label>
                            <div><input type="text" id="mem-field" name="mem" class="form-control" value="{{form.mem}}" v-model="task.mem"></div>
                        </div>                        
                    </div>
                    <div class="col-sm-4">
                        <div class="form-group">
                            <label for="disk-field" class="control-label">Disk Space(MB)</label>
                            <div><input type="text" id="disk-field" name="disk" class="form-control" value="{{form.disk}}" v-model="task.disk"></div>
                        </div>                        
                    </div>                                        
                </div>
                <div class="form-group">
                    <label for="cmd-field" class="control-label">Command</label>
                    <div><textarea type="text" id="cmd-field" name="cmd" class="form-control" v-model="task.cmd"></textarea></div>              
                </div>
                <div class="row full-bleed">
                    <div class="panel panel-default">
                        <div class="panel-heading clickable">
                            <div class="panel-title">
                                <a data-toggle="collapse" href="#DockerSetting" aria-expanded="true"><span class="si-angle-right"></span>Docker settings</a>
                                
                            </div>
                        </div>
                        <div id="DockerSetting" class="panel-collapse collapse" aria-expanded="true">
                            <div class="panel-body">
                                <div class="row">
                                    <div class="col-sm-6">
                                        <div class="form-group">
                                            <label for="image-field" class="control-label">Image</label>
                                            <div><input type="text" id="image-field" name="image" class="form-control" value="{{form.image}}" v-model="task.image"></div>
                                        </div>                                                
                                    </div>
                                    <div class="col-sm-6">
                                        <div class="form-group">
                                            <label for="network-field" class="control-label">Network</label>
                                            <div>
                                                <select type="text" id="network-field" name="network" class="form-control" v-model="task.network"> 
                                                    <option value="HOST">HOST</option>
                                                    <option value="BRIDGE">BRIDGE</option>
                                                </select></div>
                                        </div>                                                
                                    </div>                                    
                                </div>
                                <h4>Port Mappings</h4>
                                <div class="row duplicable-row">
                                    <div class="col-sm-4">
                                        <div class="form-group">
                                            <label for="portMappings[0].containerPort-field" class="control-label">Container Port</label>
                                            <div><input type="text" id="portMappings[0].containerPort-field" name="portMappings[0].containerPort" class="form-control" v-model="task.portMappings[0].containerPort"></div>
                                        </div>                                            
                                    </div>
                                    <div class="col-sm-4">
                                        <div class="form-group">
                                            <label for="portMappings[0].hostPort-field" class="control-label">Host Port</label>
                                            <div><input type="text" id="portMappings[0].hostPort-field" name="portMappings[0].hostPort" class="form-control" v-model="task.portMappings[0].hostPort"></div>                  
                                        </div>                                     
                                    </div>
                                    <div class="col-sm-4">
                                        <div class="form-group">
                                            <label for="portMappings[0].protocol-field" class="control-label">Host Port</label>
                                            <div><select type="text" id="portMappings[0].protocol-field" name="portMappings[0].protocol" class="form-control" v-model="task.portMappings[0].protocol"> 
                                                <option value="TCP">TCP</option>
                                                <option value="UDP">UDP</option>
                                            </select></div> 
                                        </div>              
                                            <div class="controls">
                                                <a href="#"><span class="si-plus"></span></a>
                                                <a href="#"><span class="si-less"></span></a>
                                            </div>                            
                                    </div>
                                </div>
                                <h4>Volumes</h4>
                                <div class="volumes">
	                                <div class="row duplicable-row volumes-row">
	                                    <div class="col-sm-4">
	                                        <div class="form-group">
	                                            <label for="volumes[0].containerPath-field" class="control-label">Container Path</label>
	                                            <div><input volumes="text" id="volumes[0].containerPath-field" name="volumes[0].containerPath" class="form-control" v-model="task.volumes[0].containerPath"></div>
	                                        </div>                                            
	                                    </div>
	                                    <div class="col-sm-4">
	                                        <div class="form-group">
	                                            <label for="volumes[0].hostPath-field" class="control-label">Host Path</label>
	                                            <div><input type="text" id="volumes[0].hostPath-field" name="volumes[0].hostPath" class="form-control" v-model="task.volumes[0].hostPath"></div>                  
	                                        </div>                                     
	                                    </div>
	                                    <div class="col-sm-4">
	                                        <div class="form-group">
	                                            <label for="volumes[0].mode-field" class="control-label">Mode</label>
	                                            <div><select type="text" id="volumes[0].mode-field" name="volumes[0].mode" class="form-control" v-model="task.volumes[0].mode"> 
	                                                <option value="RO">Read Only</option>
	                                                <option value="RW">Read and Write</option>
	                                            </select></div> 
	                                        </div>              
	                                            <div class="controls">
	                                                <a href="#"><span class="si-plus"></span></a>
	                                                <a href="#"><span class="si-less"></span></a>
	                                            </div>                            
	                                    </div>
	                                </div>   
                                </div>                             
                            </div>
                        </div>                        
                    </div>
                </div>
            </div>

            <div class="modal-footer">
                <button type="button" class="btn btn-white" data-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="submitTask()">Create</button>
            </div>
        </div>
    </div>
</div>	
</template>

<script>
	export default {
		props: {
			form: Object,
			tasks: Array
		},
		data() {
			return {
				task: {
					portMappings: [{
						hostPort: '',
						containerPort: '',
						protocol: ''
					}],
					volumes: [{
						containerPath: '',
						hostPath: '',
						mode:''
					}]
				}
			}
		},
		methods: {
			submitTask: function() {
				let _this = this				
				this.$http.post("/dist/task.json", {
					data: JSON.stringify(this.task)
				}).then(function (response) {
					// _this.tasks.push(task)
					// toastr.success('My name is Inigo Montoya. You killed my father, prepare to die!')
				}, function (response) {
					_this.tasks.push(_this.task)
					toastr.success('My name is Inigo Montoya. You killed my father, prepare to die!')
				})
			},
		}
	}
</script>