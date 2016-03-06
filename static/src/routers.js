'use strict'

import Task from './views/task.vue'

export default function(router){
    router.map({
        '/task':{				//首页
            name:'task',
            component: Task
        },    	
    })
}