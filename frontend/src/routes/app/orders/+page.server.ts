import { apiLocation } from '$lib/server/config';

declare interface OrderDTO {
    Id: string;
    Name: string;
    TaxNumber?: string;
    ZIPCode: string;
    City: string;
    Country: string;
    Address: string;
    Number: string;
    Packages: PackageDTO[];
}

declare interface PackageDTO {
    Id: string;
    Length: number;
    Width: number;
    Height: number;
    FromName: string;
    FromPhone: string;
    FromEmail: string;
    FromCountry: string;
    FromZIP: string;
    FromCity: string;
    FromAddress: string;
    FromNumber: string;
    FromOther: string;
    ToName: string;
    ToPhone: string;
    ToEmail: string;
    ToCountry: string;
    ToZIP: string;
    ToCity: string;
    ToAddress: string;
    ToNumber: string;
    ToOther: string;
}

export const load = async (event) => {
    const url = apiLocation + "/getall"
    return await event.fetch(url, {
        method: "GET",
        headers: { "Content-Type": "application/json",
            "Authorization": "Bearer " + event.cookies.get("access")
        }
    }).then(async (res) => {
        if (res.ok) {
            console

            return { orders: await res.json() as OrderDTO}
        }
    }).catch((e) => {
        console.log(e)
    })

}