<template>
    <div class="mail-box-header">

        <form method="get" action="index.html" class="pull-right mail-search">
            <div class="input-group">
                <input type="text" class="form-control input-sm" name="search" placeholder="Search all tasks...">
                <div class="input-group-btn">
                    <button type="submit" class="btn btn-sm btn-primary">
                        Search
                    </button>
                </div>
            </div>
        </form>
        <h2>
        All Tasks ({{tasks.length}})
    </h2>
        <div class="mail-tools tooltip-demo m-t-md">
            <div class="btn-group pull-right">
                <button class="btn btn-white btn-sm"><span class="si-angle-left"></span>
                </button>
                <button class="btn btn-white btn-sm"><span class="si-angle-right"></span>
                </button>

            </div>
            <button class="btn btn-white btn-sm" data-toggle="tooltip" data-placement="top" title="refresh"><span class="si-refresh-2"></span>
            </button>
            <button class="btn btn-white btn-sm" data-toggle="tooltip" data-placement="top" title="kill"><span class="si-close-circle"></span>
            </button>
            <button class="btn btn-white btn-sm" data-toggle="tooltip" data-placement="top" title="delete"><span class="si-trash"></span>
            </button>

        </div>
    </div>
    <div class="mail-box">
        <table class="table table-hover table-mail">
        	<thead>
        		<tr>
        			<th>Choose</th>
        			<th>ID</th>
        			<th>Name</th>
        			<th>Image</th>
                    <th>Slave</th>
        			<th>Cpus</th>
        			<th>Memory</th>
        			<th>Status</th>
        			<th>Time</th>
        			<th>Output</th>
        		</tr>
        	</thead>
            <tbody>
            	<tr v-for="task in tasks">
                    <td class="check-mail">
                        <input type="checkbox" class="i-checks">
                    </td>
                    <td>{{task.id}}</td>
                    <td class="mail-subject">{{task.name}}</td>
                    <td>{{task.docker_image}}</td>
                    <td>node1</td>
                    <td>{{task.cpus}}</td>
                    <td>{{task.mem}} MiB</td>
                    <td class=""><span class="label label-info">Finished</span></td>
                    <td class="mail-date">{{task.CreatedTime}}</td>
                    <td>stdout stderr</td>					
            	</tr>
            </tbody>
        </table>
    </div>	
    <sh-task-modal :form="form" :tasks="tasks"></sh-task-modal>
</template>

<script>
	import TaskModal from '../components/taskModal.vue'

	export default {
		data() {
			return {
				tasks: [],
				form: {
					cpus: 0.1,
					mem: 16,
					disk: 0,
					image: "busybox"
				}
			}
		},
		ready() {
			this.getTasks()
	        $(".i-checks").iCheck({
	            checkboxClass: "icheckbox_square-green",
	            radioClass: "iradio_square-green"
	        })	
		},
		methods:{
            getTasks (){
            	let _this = this;
			    this.$http.get('/dist/task.json').then(function (response) {
			    	if(response.data.success) {
			    		_this.$set('tasks', response.data.result)
			    	} else {
			    		//TODO 错误处理
			    	}
			    }, function (response) {
			    });
            },
		},
		components: {
			'shTaskModal': TaskModal
		}
	}
</script>