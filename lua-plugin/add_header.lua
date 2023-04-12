local library = require("library")

function envoy_on_request(request_handle)
    local headers = request_handle:headers()
    local metadata = request_handle:streamInfo():dynamicMetadata()
    metadata:set("envoy.filters.http.lua", "requestInfo", {
        requestId = headers:get("my-request-id"),
        method = headers:get(":method"),
    })
end

function envoy_on_response(response_handle)
  local requestInfoObj = response_handle:streamInfo():dynamicMetadata():get("envoy.filters.http.lua")["requestInfo"]

  local requestId = requestInfoObj.requestId
  local method = requestInfoObj.method

  if (method == 'GET') then
    if (requestId == nil or requestId == '') then
        response_handle:logInfo("Generate Request Id")
        requestId = library.RandomString()
    end
    response_handle:headers():add("my-request-id", requestId)
  end
end

