import biddrClient from '../../biddrClient/biddrClient'
import { Component } from 'react';
import '../UserSignIn/user-login';
import '../UserSignIn/user-login.scss';


export default class UserCreation extends Component{
    state = {
        username: '',
        password: '',
        displayName:''
    }
    handleUserChange = event => {
        this.setState({username: event.target.value});
    }
    handlePassChange = event => {
        this.setState({password: event.target.value});
    }
    handleDisplayChange = event => {
        this.setState({displayName: event.target.value});
    }
    
    
    handleSubmit = async event => {
        event.preventDefault();
        await biddrClient.createUser(this.state.username, this.state.displayName, this.state.password);
    }

    render(){
        return (
            <div className="login-flex">
                <div className= "logo" > User Creation</div>
                <div className="login-container">
                    <form onClick={this.handleSubmit} className="login-form">
                        <div className= "name form-field">
                        <label className="label">Username</label>
                            <input type="text" name="name" onChange={this.handleUserChange} className= "form-input"
                            placeholder="Username" required/>
                        </div>
                        
                        <div className="name form-field">
                            <label className="label">Display Name</label>
                            <input id="display-name" type="name"  onChange={this.handleDisplayChange} className="form-input" placeholder="Display Name" required/>
                        </div>

                        <div className="password form-field">
                            <label className="label" htmlFor="login-password">Password</label>
                            <input id="password" type="password"  onChange={this.handlePassChange} className="form-input" placeholder="Password" required/>
                        </div>

                        <div className="button-box" >
                            <input type="button" value="Create User" className="create-button"/>
                        </div>
                    </form>
                </div>
            </div>
        )
    }





}
