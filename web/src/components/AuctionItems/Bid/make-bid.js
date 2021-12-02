import { Component } from 'react';
import biddrClient from '../../../biddrClient/biddrClient';


export default class MakeBid extends Component{

    state={
        bid : 0.00
    }

    handleSubmit = async event => {
        event.preventDefault();
        console.log('Bid Made!')
        await biddrClient.makeBid(this.props.itemName, this.state.bid);
    }

    handleBidChange = event => {
        this.setState({bid: event.target.value});
    }

    render(){
        return(
            <div className='bid-input'>
                <label className="bid-input"><span className="bid-label">Bid</span></label>
                <input type="number" className = 'bid-amount' onChange={this.handleBidChange} placeholder='Bid Amount'/>
                <div className="bid-submit">
                    <button type="button" onClick={this.handleSubmit}>Place Bid</button>
                </div>
            </div>
        );
    };
}