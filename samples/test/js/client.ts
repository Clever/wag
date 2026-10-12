// The tests run against the default client, or against the client generated with
// x-unwrapped-fetch when WAG_UNWRAPPED_FETCH is set. Both clients have the same interface.
import * as WrappedClient from "swagger-test";
import * as UnwrappedFetchClient from "swagger-test-unwrapped-fetch";

export const unwrappedFetch = !!process.env.WAG_UNWRAPPED_FETCH;
export const Client: typeof WrappedClient = unwrappedFetch ? (UnwrappedFetchClient as any) : WrappedClient;
