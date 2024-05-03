/**
 * Copyright 2024 The Mothership Authors
 * SPDX-License-Identifier: Apache-2.0
 */

import { cn } from "src/utils"

function Skeleton({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("animate-pulse rounded-md bg-primary/10", className)}
      {...props}
    />
  )
}

export { Skeleton }
