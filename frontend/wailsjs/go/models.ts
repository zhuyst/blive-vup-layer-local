export namespace main {
	
	export class LiveConfig {
	    disable_llm: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LiveConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.disable_llm = source["disable_llm"];
	    }
	}
	export class InitRequestData {
	    code: string;
	    timestamp: number;
	    room_id: number;
	    mid: number;
	    caller: string;
	    code_sign: string;
	    config: LiveConfig;
	
	    static createFrom(source: any = {}) {
	        return new InitRequestData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.timestamp = source["timestamp"];
	        this.room_id = source["room_id"];
	        this.mid = source["mid"];
	        this.caller = source["caller"];
	        this.code_sign = source["code_sign"];
	        this.config = this.convertValues(source["config"], LiveConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Result {
	    code: number;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}

}

