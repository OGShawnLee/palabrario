export namespace models {
	
	export class UserProfile {
	    id_user: string;
	    email: string;
	    name: string;
	    user_name: string;
	    avatar_url?: string;
	    anchor_dialect: string;
	    total_xp: number;
	    current_level: number;
	    current_streak: number;
	    longest_streak: number;
	    role: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new UserProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id_user = source["id_user"];
	        this.email = source["email"];
	        this.name = source["name"];
	        this.user_name = source["user_name"];
	        this.avatar_url = source["avatar_url"];
	        this.anchor_dialect = source["anchor_dialect"];
	        this.total_xp = source["total_xp"];
	        this.current_level = source["current_level"];
	        this.current_streak = source["current_streak"];
	        this.longest_streak = source["longest_streak"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

}

