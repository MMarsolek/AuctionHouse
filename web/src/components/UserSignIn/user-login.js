import { Component } from 'react';
import './user-login.scss'
import biddrClient from '../../biddrClient/biddrClient'
import  {withCookies} from 'react-cookie'

 class UserLogIn extends Component{
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
        await biddrClient.userLogIn(this.state.name, this.state.pass);
        this.props.cookies.set("token", biddrClient.userAuth);
    }

    render(){
        if(!biddrClient.userAuth){
            biddrClient.userAuth = this.props.cookies.get("token")
        }
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
                            <input id="login-password" type="password"  onChange={this.handlePassChange} className="form-input" placeholder="Password" required/>
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

export default withCookies(UserLogIn)