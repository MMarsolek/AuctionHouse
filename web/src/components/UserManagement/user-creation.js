import biddrClient from '../../biddrClient/biddrClient'
import { Component } from 'react';

export default class UserCreation extends Component{
    state = {
        username: '',
        password: '',
        displayname:''
    }
    handleUserChange = event => {
        this.setState({username: event.target.value});
    }
    handlePassChange = event => {
        this.setState({password: event.target.value});
    }
    handleDisplayChange = event => {
        this.setState({displayname: event.target.value});
    }
    
    
    handleSubmit = async event => {
        event.preventDefault();
        await biddrClient.createUser(this.state.username, this.state.displayname, this.state.password);
    }

    render(){
        return (
            <div className="login-flex">
                <div className= "logo" > User Creation</div>
                <div className="user-creator-container">
                    <form onSubmit={this.handleSubmit} className="creation-form">
                        <div className= "form-field">
                        <label className="username-creation"><span className="hidden">Username</span></label>
                            <input type="text" name="name" onChange={this.handleUserChange} className= "form-input"
                            placeholder="Username" required/>
                        </div>
                        <div className="form form-field">
                            <label className="password-creator" htmlFor="login-password"><span className="hidden">Password</span></label>
                            <input id="password" type="password"  onChange={this.handlePassChange} className="form-input" placeholder="Password" required/>
                        </div>

                        <div className="form form-field">
                            <label className="display-name"><span className="hidden">Display Name</span></label>
                            <input id="display-name" type="name"  onChange={this.handleDisplayChange} className="form-input" placeholder="Display Name" required/>
                        </div>

                        <div className="form-field">
                            <input type="submit" value="Create User"/>
                        </div>
                    </form>
                </div>
            </div>
        )
    }





}
