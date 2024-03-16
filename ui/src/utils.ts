import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function timeToNatural(
  timestamp: string | null,
  includeInParentheses: boolean = false,
): string {
  let ret = '';
  const date = timestamp ? new Date(timestamp) : null;

  if (!date) {
    ret = '--';
  } else {
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHour = Math.floor(diffMin / 60);
    const diffDays = Math.floor(diffHour / 24);

    if (diffSec < 10) {
      ret = 'Just now';
    } else if (diffSec < 60) {
      ret = `${diffSec} seconds ago`;
    } else if (diffMin < 60) {
      ret = diffMin === 1 ? 'a minute ago' : `${diffMin} minutes ago`;
    } else if (diffHour < 24) {
      ret = diffHour === 1 ? 'an hour ago' : `${diffHour} hours ago`;
    } else if (diffDays < 10) {
      ret = diffDays === 1 ? 'a day ago' : `${diffDays} days ago`;
    } else {
      ret = date.toISOString().replace('T', ' ').replace('Z', '') + ' UTC'; // Custom format as per requirement
      // Return YYYY-MM-DD HH:MM
      ret = ret.slice(0, 16);
    }
  }

  if (includeInParentheses && ret !== '--') {
    ret = ` (${ret})`;
  }

  return ret;
}

export function capitalizeFirstLetter(str: string) {
  return (str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()).replace(
    '_',
    ' ',
  );
}
