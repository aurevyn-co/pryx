import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import SuperadminDashboard from './SuperadminDashboard';

describe('SuperadminDashboard', () => {
    const mockOnLogout = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        
        // Setup default mock responses
        vi.spyOn(global, 'fetch').mockResolvedValue({
            ok: true,
            json: async () => ({
                totalUsers: 1247,
                activeUsers: 892,
                newUsersToday: 23,
                totalDevices: 3421,
                onlineDevices: 2187,
                offlineDevices: 1234,
                totalSessions: 15432,
                totalCost: 2847.50,
                avgCostPerUser: 2.28,
            }),
        } as Response);
    });

    it('renders dashboard header with title', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        expect(screen.getByText('Pryx Superadmin Dashboard')).toBeInTheDocument();
    });

    it('renders navigation tabs', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        expect(screen.getByText('Overview')).toBeInTheDocument();
        expect(screen.getByText('Users')).toBeInTheDocument();
        expect(screen.getByText('Devices')).toBeInTheDocument();
        expect(screen.getByText('Costs')).toBeInTheDocument();
        expect(screen.getByText('System Health')).toBeInTheDocument();
    });

    it('displays logout button', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        expect(screen.getByText('Logout')).toBeInTheDocument();
    });

    it('calls onLogout when logout button is clicked', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        fireEvent.click(screen.getByText('Logout'));
        expect(mockOnLogout).toHaveBeenCalledTimes(1);
    });

    it('switches to users tab when clicked', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        fireEvent.click(screen.getByText('Users'));
        expect(screen.getByText('User Management')).toBeInTheDocument();
    });

    it('switches to devices tab when clicked', () => {
        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        fireEvent.click(screen.getByText('Devices'));
        expect(screen.getByText('Device Fleet')).toBeInTheDocument();
    });

    it('switches to system health tab when clicked', () => {
        // Setup health mock
        vi.spyOn(global, 'fetch').mockResolvedValue({
            ok: true,
            json: async () => ({
                runtimeStatus: 'healthy',
                apiLatency: 45,
                errorRate: 0.001,
                dbStatus: 'connected',
                queueDepth: 12,
                activeConnections: 456,
            }),
        } as Response);

        render(<SuperadminDashboard onLogout={mockOnLogout} />);
        fireEvent.click(screen.getByText('System Health'));
        expect(screen.getByText('HEALTHY')).toBeInTheDocument();
    });
});
