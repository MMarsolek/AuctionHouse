import { Component } from 'react';
import  { withCookies } from 'react-cookie'
import { Navigate } from 'react-router-dom';
import { AuthContext } from '../../auth-provider';
import { withNotifications } from '../../utils';
import './user-login.scss'

 class UserLogIn extends Component{
    static contextType = AuthContext;
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
    
    handleSubmit = async event => {
        event.preventDefault();
        const loginSuccessful = await this.context.login(this.state.name, this.state.pass)
        if (!loginSuccessful) {
            this.props.notify('Login not found or password was incorrect', 'error');
        }
    }

    render(){
        return (
            <div className="login-flex">
                { this.context.user && <Navigate to="/auctions" replace={true} />}
                <div className= "logo" >AuctionHouse Log In</div>
                <div className="login-container">
                    <form onClick={this.handleSubmit} className="login-form">
                        <div className= "form-field">
                        <label className="username"><span className="hidden">Username</span></label>
                            <input type="text" name="name" onChange={this.handleUserChange} className= "form-input"
                            placeholder="Username" required/>
                        </div>
                        <div className="form-field">
                            <label className="lock" htmlFor="login-password"><span className="hidden">Password</span></label>
                            <input id="login-password" type="password"  onChange={this.handlePassChange} className="form-input" placeholder="Password" required/>
                        </div>

                        <div className="button-box">
                            <input type="button" value="Log in"className="create-button"/>
                        </div>
                    </form>
                </div>
            </div>
        )
    }
}

export default withNotifications(withCookies(UserLogIn));