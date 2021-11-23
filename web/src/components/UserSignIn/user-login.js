import axios from 'axios'
import { Component } from 'react';
import './user-login.css'

export default class UserLogIn extends Component{
    state = {
        name: '',
        pass: '',
    }

    handleUserChange = event => {
        this.setState({name: event.target.value});
    }
    handlePassChange = event => {
        this.setState({pass: event.target.value});
    }
    
    
    handleSubmit = event => {
        event.preventDefault();


        axios.post('http://localhost:8080/api/v1/users/login', {
            username: this.state.name,
            password: this.state.pass
        })
            .then(res => {
                console.log(res);
                console.log(res.data);
            })
    }

    render(){
        return (
            <form onSubmit={this.handleSubmit}>
                <label>
                    User Name:
                    <input type="text" name="name" onChange={this.handleUserChange}/>
                </label>
                <label>
                    Password:
                    <input type="password" name="pass" onChange={this.handlePassChange}/>
                </label>
                <button type="submit" className='login-button'> Log In
                </button>
            </form>
        )
    }
}