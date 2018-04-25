class Index extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleSend = this.handleSend.bind(this);
    this.authenticate = this.authenticate.bind(this);
    this.state = {username: '', password: '', authenticated: false};
  }
  componentDidMount(){
   this.authenticate();
  }
  authenticate(){
    axios.get('/authenticate')
      .then(res => {
        if(res.data == 'Authenticated'){
          this.setState({authenticated: true});
        }
      })
      .catch(error => {
        console.log(error)
      });

  }
  handleChange(event) {
    this.setState({[event.target.name]: event.target.value}); //Dynamic keys https://stackoverflow.com/questions/29280445/reactjs-setstate-with-a-dynamic-key-name
  }
  handleSend(event) {
    axios.post('/' + event.target.name, {
        "username": this.state.username,
        "password": this.state.password,
      })
      .then(res => {
        this.authenticate();
      })
      .catch(error => {
        console.log(error)
      });
    event.preventDefault();
    
  }
  render() {
    if (this.state.authenticated){
      return <p> Hello </p>
    }
    return (
      <form>
      username:
      <input type="text" name="username" value={this.state.username} onChange={this.handleChange} />
      password:
      <input type="text" name="password" value={this.state.password} onChange={this.handleChange} />
      <button name="login" onClick={this.handleSend}>Login</button>
      <button name="register" onClick={this.handleSend}>Register</button>
      </form>
    );
  }
}

ReactDOM.render( <Index/>, document.querySelector("#root"));
