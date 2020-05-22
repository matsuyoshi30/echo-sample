const app = new Vue({
  el: '#app',
  data: {
    todos: [],
    taskname: '',
  },
  created() { this.update() },
  methods: {
    add: () => {
      const payload = {'taskname': app.taskname}
      axios.post('/api/todos', payload)
           .then(() => {
             app.taskname = ''
             app.update()
           })
           .catch((err) => {
             alert(err.response.data.error)
           });
    },
    check: (todo) => {
      let params = new URLSearchParams();
      params.append('completed', !todo.completed);
      axios.put('/api/todos/' + todo.id, params)
           .then(() => {
             app.update()
           })
           .catch((err) => {
             alert(err.response.data.error)
           });
    },
    remove: (todo) => {
      axios.delete('/api/todos/' + todo.id)
           .then(() => {
             app.update()
           })
           .catch((err) => {
             alert(err.response.data.error)
           });
    },
    update: () => {
      axios.get('/api/todos')
           .then((response) => app.todos = response.data || [])
           .catch((error) => console.log(error));
    },
  }
})
