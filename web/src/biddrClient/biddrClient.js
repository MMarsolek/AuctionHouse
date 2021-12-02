import axios from 'axios';

class BiddrClient{

    userAuth = '';
    domain = '';

    constructor(baseUrl){
        this.domain = baseUrl
    }
    
    async requestConfig(url, method, {headers, body} = {}){
        const myParameters = {
            url : url, 
            method : method, 
            baseURL : this.domain,  
            validateStatus :  status => status < 500,
            headers: {
                ...(headers || {}),
                authorization: `Bearer ${this.userAuth}`,
            }
        };

        if(body){
            myParameters['data'] = body;
        };
        const myResults = await axios(myParameters);
        
        return myResults.data;
    }

    async createUser(user, display, pass){
        return await this.requestConfig('/api/v1/users/', 'post', {
            body: {
            username: user,
            password: pass,
            displayName: display
        }})
    }

    async getUserInformation(userName){
        const encodedUser = encodeURIComponent(userName)
        const response = await  this.requestConfig(`/api/v1/users/${encodedUser}`, 'get');
        return response;
    }

    async userLogIn(username, password){
        const response = await this.requestConfig('/api/v1/users/login', 'post',{
            body: {
            username: username,
            password: password,
            
        }});
        this.userAuth = response['authToken'];
        return response;
    }

    async userLogout() {
        this.authToken = '';
    }

    returnHeaders(){
        return ({'authorization':'Bearer ' + this.userAuth})
    }

    async getHighestBidForAll(){
        return await this.requestConfig('/api/v1/auctions/bids', 'get');
     }

    async getHighestBidForOneItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'get');
    }

    async makeBid(name, bid){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'post', {
            body: {bidAmount: bid},
        });
    }
    
    async getAllItems(){
        return await this.requestConfig(`/api/v1/auctions/items`, 'get')
    }

    async createNewItem(description, image, name){
        console.log(name);
        return await this.requestConfig('/api/v1/auctions/items', 'post', {
            body:
            {   description,
                image,
                name
            },
        });
    }

    async deleteItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'delete');
    }

    async getSpecificItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'get');
    }

    async updateItem(description, image, name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'put', {
            body: {
                description: description,
                image: image,
                name: name
            }
        });
    }

    async establishWebSocket(){
        return await this.requestConfig('/api/v1/ws', 'get');
    }

}

export default new BiddrClient("http://localhost:8080")
