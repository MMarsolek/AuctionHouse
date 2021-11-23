import axios from 'axios'
import react, { Component } from 'react';

export default class UserLogIn extends Component{
    state = {
        userInfo: []
    }
    constructor(){
        super();
        axios.get('http://localhost:1234/api/v1/api/v1/users/login')
            .then(res => {
                console.log(res.data)
                this.setState({ userInfo: res.data});
        })
    }

    render(){
        return (
            <ul>
                {this.state.userInfo.indexOf(0)}
            </ul>
        )
    }
}