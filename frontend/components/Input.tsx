import { IconProp } from "@fortawesome/fontawesome-svg-core"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import classNames from "classnames"
import React from "react"

type InputProps = {
  icon?: IconProp
  error?: string
} & React.ComponentProps<"input">

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ icon, className, error, ...props }, ref) => {
    return (
      <div className={className}>
        <div className="relative">
          {icon && (
            <FontAwesomeIcon
              icon={icon}
              className="h-5 w-5 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400"
            />
          )}

          <input
            {...props}
            className={classNames(
              "border rounded w-full min-h-30 py-2 pr-2.5 pl-9 focus:ring-2 focus:outline-none transition-shadow focus:ring-primary-light",
              error ? "border-red-700" : "border-gray-300"
            )}
            ref={ref}
          />
        </div>
        {error ? (
          <div className="text-sm text-red-700 mt-0.5 mx-2.5">{error}</div>
        ) : (
          ""
        )}
      </div>
    )
  }
)

export default Input
