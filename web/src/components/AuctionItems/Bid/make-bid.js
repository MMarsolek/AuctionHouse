import { Component } from 'react';
import { withNotifications } from '../../../utils';
import biddrClient from '../../../biddrClient/biddrClient';


class MakeBid extends Component{

    state = {
        bid : 0.00
    }

    handleSubmit = async event => {
        event.preventDefault();
        try {
            await biddrClient.makeBid(this.props.itemName, this.state.bid);
            this.props.updateItem({ bidAmount: this.state.bid});
        } catch (ex) {
            this.props.notify(`Unable to make a bid: ${ex.message}`, 'error');
        }
    }

    handleBidChange = event => {
        this.setState({bid: parseInt(event.target.value)});
    }

    render(){
        return(
            <div className='bid-input'>
                <label className="bid-input"><span className="bid-label">Bid</span></label>
                <input type="number" className='bid-placement' onChange={this.handleBidChange} placeholder='Bid Amount'/>
                <div className="bid-submit">
                    <button type="button" onClick={this.handleSubmit}>Place Bid</button>
                </div>
            </div>
        );
    };
}

export default withNotifications(MakeBid);