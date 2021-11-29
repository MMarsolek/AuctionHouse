import { Component } from 'react';
import './user-login.scss'
import userLogIn from '../../biddrClient/biddrClient'

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
        userLogIn(this.name, this.pass)
    }

    render(){
        return (
            <div className="login-flex">
                <div className= "logo" > AuctionHouse Log In</div>
                <div className="login-container">
                    <form onSubmit={this.handleSubmit} className="login-form">
                        <div className= "form-field">
                        <label className="username"><span className="hidden">Username</span></label>
                            <input type="text" name="name" onChange={this.handleUserChange} className= "form-input"
                            placeholder="Username" required/>
                        </div>

                        <div className="form form-field">
                            <label className="lock" htmlFor="login-password"><span className="hidden">Password</span></label>
                            <input id="login-password" type="password" className="form-input" placeholder="Password" required/>
                        </div>

                        <div className="form-field">
                            <input type="submit" value="Log in"/>
                        </div>
                    </form>
                </div>
            </div>

        )
    }
}