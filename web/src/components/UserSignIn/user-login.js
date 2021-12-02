import { Component } from 'react';
import  { withCookies } from 'react-cookie'
import { Navigate } from 'react-router-dom';
import { AuthContext } from '../../auth-provider';
import './user-login.scss'

 class UserLogIn extends Component{
    static contextType = AuthContext;
    state = {
        name: '',
        pass: '',
        proceed: false,
    }


    handleUserChange = event => {
        this.setState({name: event.target.value});
    }
    handlePassChange = event => {
        this.setState({pass: event.target.value});
    }
    
    handleSubmit = async event => {
        event.preventDefault();
        if (await this.context.login(this.state.name, this.state.pass)) {
            this.setState({proceed: true});
        }
    }

    async componentDidMount() {
        const wasSuccessful = await this.context.refreshSession();
        if (wasSuccessful) {
            this.setState({proceed: true});
        }
    }

    render(){
        return (
            <div className="login-flex">
                { this.state.proceed && <Navigate to="/auctions" replace={true} />}
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