import { Component } from 'react';
import biddrClient from '../../../biddrClient/biddrClient';


export default class MakeBid extends Component{

    state={
        bid : 0,
        itemName : ''
    }

    handleClick = async event => {

        event.preventDefault();
        this.setState({itemName : this.props.MakeBid});
        await biddrClient.userLogIn(this.state.itemName, this.state.bid);
    }

    handleBidChange = event => {
        this.setState({bid: event.target.value});
    }

    render(){
        return(
            <div className='bid-maker'>
                <input type="number" className = 'bid-amount' onChange={this.setState.bid}>
                    Bid Amount
                </input>
                <button className= 'bid-button' onClick={this.handleClick}> Make Bid

                </button>
            </div>    
        
        );
    };
}