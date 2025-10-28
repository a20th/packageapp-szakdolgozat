// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			loggedIn: Boolean
			access: string | null
			refresh: string | null
			profile: ProfileData | null
		}
		interface ProfileData{
			name: string
			phone: string
			email: string
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
