import { Component } from 'react';
import biddrClient from '../../../biddrClient/biddrClient';


export default class MakeBid extends Component{

    state={
        bid : 0,
    }

    handleSubmit = async event => {

        event.preventDefault();
        await biddrClient.makeBid(this.props.itemName, this.state.bid);
    }

    handleBidChange = event => {
        this.setState({bid: event.target.value});
    }

    render(){
        return(
            <form onSubmit={this.handleSubmit} className="bid-form">
                <div className='bid-maker'>
                    <div className= 'bid-input'>
                    <label className="bid-input"><span className="bid-label">Bid</span></label>
                    <input type="number" className = 'bid-amount' onChange={this.handleBidChange} placeholder='Bid Amount'/>
                    </div>
                    <div className="bid-submit">
                        <input type="submit" value="Place Bid" className='bid-button'/>
                    </div>
                </div>    
            </form>
        );
    };
}