export namespace params {
	
	export class Message {
	    id: number;
	    hash: string;
	    content: string;
	    name: string;
	    timestamp: number;
	    isStored: boolean;
	    wakuTimestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.hash = source["hash"];
	        this.content = source["content"];
	        this.name = source["name"];
	        this.timestamp = source["timestamp"];
	        this.isStored = source["isStored"];
	        this.wakuTimestamp = source["wakuTimestamp"];
	    }
	}

}

