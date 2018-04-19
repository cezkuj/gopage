class Index extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleLogin = this.handleLogin.bind(this);
    this.handleRegister = this.handleRegister.bind(this);
    this.state = {username: '', password: ''};
  }
  handleChange(event) {
    this.setState({[event.target.name]: event.target.value}); //Dynamic keys https://stackoverflow.com/questions/29280445/reactjs-setstate-with-a-dynamic-key-name
  }
  handleLogin(event) {
    fetch('/login', {
      method: 'POST',
      body: "username=" + this.state.username + "&password=" + this.state.password,
    });
    event.preventDefault();
    
  }

  handleRegister(event) {
    fetch('/register', {
      method: 'POST',
      body: "username=" + this.state.username + "&password=" + this.state.password,
    });
    event.preventDefault();
  }
  render() {
    return (
     <form onSubmit={this.handleLogin}>
      username:
      <input type="text" name="username" value={this.state.username} onChange={this.handleChange} />
      password:
      <input type="text" name="password" value={this.state.password} onChange={this.handleChange} />
      <button type="submit">Login</button>
      <button onClick={this.handleRegister}>Register</button>
    </form>
    );
  }
}

ReactDOM.render( <Index/>, document.querySelector("#root"));

