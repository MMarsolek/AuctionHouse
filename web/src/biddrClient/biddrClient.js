import axios from 'axios'

class BiddrClient{

    userAuth = '';

    constructor(baseUrl){
        this.domain = baseUrl
    }

    createUser(user, display, pass){
        return axios.post(`${this.domain}/api/v1/users/`, { 
            username: user,
            password: pass,
            displayName: display
        },{
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        });
    }

    getUserInformation(userName){
        const encodedUser = encodeURIComponent(userName)
        return axios.get(`${this.domain}/api/v1/users/${encodedUser}`, {
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            this.userAuth = res.data['authToken'];
            return res.data
        });
    }

    userLogIn(user, pass){
        return axios.post(`${this.domain}/api/v1/users/login`, {   
            username: user,
            password: pass
        },{ validateStatus: status => {
                return status < 500;
            }}).then(res => {
            return res.data
        })
    }

    getHighestBidForAll(){
        return axios.get(`${this.domain}/api/v1/auctions/bids`,  { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    getHighestBidForOneItem(name){
        const encodeditem = encodeURIComponent(name)
        return axios.get(`${this.domain}/api/v1/auctions/bids${encodeditem}`,  { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    makeBid(name, bid){
        const encodeditem = encodeURIComponent(name)
        return axios.post(`${this.domain}/api/v1/auctions/bids${encodeditem}`, {
            'bidAmount': bid}, { 
            headers: {
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    getAllItems(){
        return axios.get(`${this.domain}/api/v1/auctions/items`,{ 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    createNewItem(description, image, name){
        const encodeditem = encodeURIComponent(name)
        return axios.post(`${this.domain}/api/v1/auctions/bids${encodeditem}`,  {
            description: description,
            image: image,
            name: name
        }, { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    deleteItem(name){
        const encodeditem = encodeURIComponent(name)
        return axios.delete(`${this.domain}/api/v1/auctions/bids${encodeditem}`, { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    getSpecificItem(name){
        const encodeditem = encodeURIComponent(name)
        return axios.get(`${this.domain}/api/v1/auctions/bids${encodeditem}`, { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    updateItem(description, image, name){
        const encodeditem = encodeURIComponent(name)
        return axios.put(`${this.domain}/api/v1/auctions/bids${encodeditem}`,  {
            description: description,
            image: image,
            name: name
        }, { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

    establisWebSocket(){
        return axios.get(`${this.domain}/api/v1/ws`, { 
            headers:{
                authenication: 'Bearer ' + this.userAuth
            }  
        }, { 
            validateStatus: status => {
                return status < 500;
            }
        }).then(res => {
            return res.data
        })
    }

}
export default new BiddrClient("http://localhost:8080")