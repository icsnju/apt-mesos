<template>
        <div class="row">
            <div class="col-sm-2">
                <div class="ibox float-e-margins">
                    <div class="ibox-content mailbox-content">
                        <div class="file-manager">
                            <button class="btn btn-block btn-primary compose-mail" data-toggle="modal" data-target="#addTask">Add Task</button>
                            <div class="space-25"></div>
							<sh-task-filter></sh-task-filter>

                            <h5 class="tag-title">Label</h5>
                            <ul class="tag-list" style="padding: 0">
                            </ul>
                            <div class="clearfix"></div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-sm-10 animated fadeInRight">
				<sh-task-list>
				</sh-task-list>
            </div>
        </div>
</template>

<script>
import Filter from '../components/taskFilter.vue'
import List from '../components/taskList.vue'
export default {

	components: {
		shTaskFilter: Filter,
		shTaskList: List
	}
}
</script>

